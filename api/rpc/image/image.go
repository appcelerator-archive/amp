package image

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"docker.io/go-docker/api/types"
	"github.com/appcelerator/amp/pkg/cloud"
	"github.com/appcelerator/amp/pkg/cloud/aws"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/docker/distribution/reference"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server server information
type Server struct {
	Docker   *docker.Docker
	Provider cloud.Provider
}

type loadStatus struct {
	Status string `json:"status"`
	ID     string `json:"id"`
	Stream string `json:"stream"`
}
type imageList struct {
	Repositories []string `json:"repositories"`
}
type tagList struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
type removeError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type removeContent struct {
	Errors []removeError `json:"errors"`
}

func registryEndpoint(ctx context.Context, provider cloud.Provider) (string, error) {
	var registryEndpoint string
	var err error
	switch provider {
	case cloud.ProviderAWS:
		registryEndpoint, err = aws.InternalRegistry(ctx)
		if err != nil {
			return "", err
		}
	case cloud.ProviderLocal:
		return "", errors.New("no registry for local deployment")
	default:
		return "", errors.New("not yet implemented for this provider")
	}
	return strings.ToLower(registryEndpoint), nil
}

func tagImage(image, prefix string) (string, error) {
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return "", errors.New("not a valid repository/tag")
	}
	if _, isCanonical := ref.(reference.Canonical); isCanonical {
		return "", errors.New("refusing to create a tag with a digest reference")
	}
	taggedImage := image
	hostname := reference.Domain(ref)
	if hostname != prefix {
		taggedImage = prefix + "/" + image
	}
	return taggedImage, nil
}

func (s *Server) ImagePush(ctx context.Context, in *PushRequest) (*PushReply, error) {
	log.Infoln("[image] Push", in.Name)
	reg, err := registryEndpoint(ctx, s.Provider)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "registry endpoint: %s", err.Error())
	}
	ac := types.AuthConfig{Username: "none"}
	jsonString, err := json.Marshal(ac)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal authconfig")
	}
	dst := make([]byte, base64.URLEncoding.EncodedLen(len(jsonString)))
	base64.URLEncoding.Encode(dst, jsonString)
	authConfig := string(dst)
	imagePushOptions := types.ImagePushOptions{RegistryAuth: authConfig}

	taggedImage, err := tagImage(in.Name, reg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	// load image from archive
	loadResp, err := s.Docker.GetClient().ImageLoad(ctx, bytes.NewReader([]byte(in.Data)), false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "image load: %s", err.Error())
	}
	log.Infoln("Image loaded")
	defer loadResp.Body.Close()
	var tmpImage string
	if loadResp.Body != nil && loadResp.JSON {
		lResp, err := ioutil.ReadAll(loadResp.Body)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to read the ImageLoad response: %s", err.Error())
		}
		re := regexp.MustCompile(`"stream":"Loaded image: (.*)\\n`)
		match := re.FindStringSubmatch(string(lResp))
		if len(match) != 2 {
			return nil, status.Errorf(codes.Internal, "failed to match the ImageLoad response: %s", lResp)
		}
		tmpImage = match[1]
	} else {
		return nil, status.Errorf(codes.Internal, "image load response is empty")
	}
	// tag image
	if err = s.Docker.GetClient().ImageTag(ctx, tmpImage, taggedImage); err != nil {
		return nil, status.Errorf(codes.Internal, "trying to tag [%s] to [%s]: %s", tmpImage, taggedImage, err.Error())
	}

	pushResp, err := s.Docker.GetClient().ImagePush(ctx, taggedImage, imagePushOptions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "image push: %s", err.Error())
	}
	defer pushResp.Close()
	body, err := ioutil.ReadAll(pushResp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "image push response: %s", err.Error())
	}
	re := regexp.MustCompile(`digest: (sha256:[^\s]*) size:`)
	match := re.FindStringSubmatch(string(body))
	if len(match) != 2 {
		log.Infoln("Image push failed:", string(body))
		return nil, status.Errorf(codes.Internal, string(body))
	}
	// remove the images
	if _, err := s.Docker.GetClient().ImageRemove(ctx, tmpImage, types.ImageRemoveOptions{}); err != nil {
		log.Infoln("Failed to remove image", tmpImage)
	}
	if _, err := s.Docker.GetClient().ImageRemove(ctx, taggedImage, types.ImageRemoveOptions{}); err != nil {
		log.Infoln("Failed to remove image", taggedImage)
	}
	log.Infoln("Successfully pushed image:", in.Name, "digest:", match[1])
	return &PushReply{Digest: match[1]}, nil
}

func (s *Server) ImageList(ctx context.Context, in *ListRequest) (*ListReply, error) {
	protocol := "http"
	var err error
	log.Infoln("[image] List")
	reg, err := registryEndpoint(ctx, s.Provider)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	hclient := &http.Client{}

	resp, err := hclient.Get(protocol + "://" + reg + "/v2/_catalog")
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	defer resp.Body.Close()
	repositories, err := ioutil.ReadAll(resp.Body)
	//log.Infoln(string(repositories))
	var il imageList
	if err := json.Unmarshal(repositories, &il); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	reply := &ListReply{}
	for _, repo := range il.Repositories {
		r, err := hclient.Get(protocol + "://" + reg + "/v2/" + repo + "/tags/list")
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		tags, err := ioutil.ReadAll(r.Body)
		var tl tagList
		if err := json.Unmarshal(tags, &tl); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		imageEntries := []*ImageEntry{}
		for _, tag := range tl.Tags {
			req, err := http.NewRequest("HEAD", protocol+"://"+reg+"/v2/"+repo+"/manifests/"+tag, nil)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
			m, err := hclient.Do(req)
			digest := ""
			if err != nil {
				log.Infoln("Failed to fetch the digest", err.Error())
			} else {
				digest = m.Header.Get("Docker-Content-Digest")
			}
			imageEntries = append(imageEntries, &ImageEntry{Tag: tag, Digest: digest})
		}
		reply.Entries = append(reply.Entries, &RepositoryEntry{Name: tl.Name, Entries: imageEntries})
		r.Body.Close()
	}
	log.Infoln("Successfully listed images")
	return reply, nil
}

func (s *Server) ImageRemove(ctx context.Context, in *RemoveRequest) (*RemoveReply, error) {
	protocol := "http"
	log.Infoln("[image] Remove", in.String())
	reg, err := registryEndpoint(ctx, s.Provider)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	hclient := &http.Client{}
	req, err := http.NewRequest("DELETE", protocol+"://"+reg+"/v2/"+in.Name+"/manifests/"+in.Digest, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp, err := hclient.Do(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete image: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	log.Infoln(string(body))
	if resp.StatusCode != http.StatusAccepted {
		var removeErr removeContent
		if err := json.Unmarshal(body, &removeErr); err != nil {
			return nil, status.Errorf(codes.Internal, "image removal failed, also: %s", err.Error())
		}
		var reply string
		for _, e := range removeErr.Errors {
			reply = fmt.Sprintf("%s %s: %s", reply, e.Code, e.Message)
		}
		return nil, status.Errorf(codes.Internal, reply)
	}
	log.Infoln("Successfully removed image")
	return &RemoveReply{}, nil
}
