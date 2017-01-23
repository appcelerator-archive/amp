package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var (
	verbose         bool
	protocInstalled bool
	dockerArgs      []string
	protocArgs      []string
)

func init() {
	flag.BoolVar(&verbose, "v", false, "print matched proto files")
	flag.BoolVar(&protocInstalled, "protoc", false, "if true, use system protoc, else run in a new container")

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dockerArgs = []string{
		"run", "-t", "--rm",
		"--name", "protoc",
		"-u", fmt.Sprintf("%s:%s", strconv.Itoa(os.Getuid()), strconv.Itoa(os.Getgid())),
		"-v", fmt.Sprintf("%s:%s", wd, "/go/src/github.com/appcelerator/amp"),
		"appcelerator/gotools", "protoc",
	}

	protocArgs = []string{
		"-I", "/go/src/",
		"-I", "/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis",
		"--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:/go/src/",
		"--grpc-gateway_out=logtostderr=true:/go/src",
		"--swagger_out=logtostderr=true:/go/src/",
	}
}

func main() {
	flag.Parse()

	// find and compile all *.proto files not in excluded dirs
	filepath.Walk(".", func(p string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if p == "vendor" || p == ".git" || p == ".glide" {
			return filepath.SkipDir
		}
		if filepath.Ext(p) == ".proto" {
			protoc(path.Join("/go/src/github.com/appcelerator/amp/", p))
		}
		return nil
	})
}

func protoc(p string) {
	if verbose {
		fmt.Println(p)
	}

	var cmd string
	var args []string
	if protocInstalled {
		cmd = "protoc"
		args = append(protocArgs, p)
	} else {
		cmd = "docker"
		args = append(append(dockerArgs, protocArgs...), p)
	}

	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		fmt.Println(args)
		fmt.Println(string(out))
		panic(err)
	}
}
