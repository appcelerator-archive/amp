package object_store

import (
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/object_stores"
	"github.com/appcelerator/amp/pkg/cloud"
	"github.com/appcelerator/amp/pkg/cloud/aws"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement objectStore.ObjectStoreServer
type Server struct {
	Accounts     accounts.Interface
	ObjectStores object_stores.Interface
	Provider     cloud.Provider
	Region       string
}

func convertError(err error) error {
	switch err {
	case object_stores.InvalidName:
		return status.Errorf(codes.InvalidArgument, err.Error())
	case object_stores.AlreadyExists, object_stores.AlreadyOwnedByYou:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case object_stores.NotFound:
		return status.Errorf(codes.NotFound, err.Error())
	case accounts.NotAuthorized:
		return status.Errorf(codes.PermissionDenied, err.Error())
	case object_stores.NotImplemented:
		return status.Errorf(codes.Unimplemented, err.Error())
	}
	return status.Errorf(codes.Internal, err.Error())
}

// Create implements objectStore.Server
func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	log.Infoln("[objectstore] Create,", in.String())

	// validate cloud context
	switch s.Provider {
	case cloud.ProviderAWS:
	default:
		return nil, convertError(object_stores.NotImplemented)
	}

	// Check if object store already exists in AMP
	ostore, err := s.ObjectStores.GetByName(ctx, in.Name)
	if err != nil {
		return nil, convertError(err)
	}

	if ostore == nil {
		if ostore, err = s.ObjectStores.Create(ctx, in.Name); err != nil {
			return nil, convertError(err)
		}
	} else {
		return nil, convertError(object_stores.AlreadyOwnedByYou)
	}

	switch s.Provider {
	case cloud.ProviderAWS:
		location, err := aws.S3CreateBucket(in.Name, in.Acl, s.Region)
		if err != nil {
			log.Infoln(err)
			switch err.Error() {
			case aws.ErrMsgBucketAlreadyOwnedByYou:
				err = object_stores.AlreadyOwnedByYou
			case aws.ErrMsgBucketAlreadyExists:
				err = object_stores.AlreadyExists
			default:
			}
			s.ObjectStores.Delete(ctx, ostore.Id)
			return nil, convertError(err)
		}
		if err = s.ObjectStores.UpdateLocation(ctx, ostore.Id, location); err != nil {
			return nil, convertError(err)
		}
		ostore.Location = location
		log.Infoln("new object store successfully created:", location)
		// TODO: add IAM policy to instance profile
	default:
	}

	log.Infoln("[objectstore] Success: created object store")
	return &CreateReply{
		Id:       ostore.Id,
		Name:     in.Name,
		Location: ostore.Location,
	}, nil
}

// List implements object.Server
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	var buckets []*ObjectStoreEntry
	log.Infoln("[objectstore] List", in.String())

	// first check what is has been stored by AMP
	ostores, err := s.ObjectStores.List(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	// then check with the reality
	switch s.Provider {
	case cloud.ProviderAWS:
		result, err := aws.S3ListBuckets(s.Region)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		log.Infof("[objectstore] AMP registered %d object stores, AWS API returned %d S3 buckets", len(ostores), len(*result))
		for _, ostore := range ostores {
			bucketExists := false
			for _, bucket := range *result {
				if bucket == ostore.Name {
					l, err := aws.S3GetBucketLocation(ostore.Name, s.Region)
					if err != nil {
						log.Infoln("failed to get location for s3 bucket", ostore.Name)
						return nil, status.Errorf(codes.Internal, err.Error())
					}
					buckets = append(buckets, &ObjectStoreEntry{
						ObjectStore: ostore,
						Region:      l,
					})
					bucketExists = true
					break
				}
			}
			if !bucketExists {
				log.Infof("Warning: can't find bucket %s in s3", ostore.Name)
				buckets = append(buckets, &ObjectStoreEntry{
					ObjectStore: ostore,
					Missing:     true,
				})
			}
		}
	default:
	}
	if len(ostores) != len(buckets) {
		log.Infoln(len(ostores), "object stores registered in AMP, but", len(buckets), "found")
	}
	log.Infof("[objectstore] Success: list returned %d object stores", len(buckets))
	return &ListReply{Entries: buckets}, nil
}

// Remove implements object.Server
func (s *Server) Remove(ctx context.Context, in *RemoveRequest) (*RemoveReply, error) {
	log.Infoln("[objectstore] Remove", in.String())

	// Retrieve the object store
	ostore, err := s.ObjectStores.GetByFragmentOrName(ctx, in.Name)
	if err != nil {
		return nil, convertError(err)
	}
	if ostore == nil {
		return nil, object_stores.NotFound
	}

	switch s.Provider {
	case cloud.ProviderAWS:
		if err := aws.S3DeleteBucket(ostore.Name, s.Region, in.Force); err != nil {
			switch err.Error() {
			case aws.ErrMsgNoSuchBucket:
				log.Infof("Bucket %s does not exist. Will ignore this error and delete it from AMP", ostore.Name)
			default:
				log.Infof("AWS error: %v", err)
				return nil, status.Errorf(codes.Internal, err.Error())
			}
		}
	default:
	}

	if err = s.ObjectStores.Delete(ctx, ostore.Id); err != nil {
		return nil, convertError(err)
	}

	log.Infoln("[objectstore] Success: removed", ostore.Name)
	return &RemoveReply{ostore.Name}, nil
}

// Forget implements object.Server
func (s *Server) Forget(ctx context.Context, in *ForgetRequest) (*ForgetReply, error) {
	log.Infoln("[objectstore] Forget", in.String())

	// Retrieve the object store
	ostore, err := s.ObjectStores.GetByFragmentOrName(ctx, in.Name)
	if err != nil {
		return nil, convertError(err)
	}
	if ostore == nil {
		return nil, object_stores.NotFound
	}

	if err = s.ObjectStores.Delete(ctx, ostore.Id); err != nil {
		return nil, convertError(err)
	}

	log.Infoln("[objectstore] Success: forget", ostore.Name)
	return &ForgetReply{Name: ostore.Name}, nil
}
