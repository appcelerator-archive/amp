package main

import "path/filepath"
import "fmt"
import "os"
import "os/exec"
import "strconv"
import "flag"

var protoc bool

func init() {
	flag.BoolVar(&protoc, "protoc", false, "protoc available, otherwise try docker")
}

func main() {
	flag.Parse()
	matches := []string{}
	paths := []string{}
	filepath.Walk(".", func(path string, info os.FileInfo, walkErr error) (err error) {
		if walkErr != nil {
			return walkErr
		}
		if path == "vendor" || path == ".git" || path == ".glide" {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".proto" {
			matches = append(matches, path)
		}
		return
	})
	for _, match := range matches {
		paths = append(paths, "/go/src/github.com/appcelerator/amp/"+match)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dockerArgs := []string{
		"run",
		"-u",
		strconv.Itoa(os.Getgid()) + ":" + strconv.Itoa(os.Getgid()),
		"--rm",
		"--name",
		"protoc",
		"-t",
		"-v",
		wd + ":/go/src/github.com/appcelerator/amp",
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		"appcelerator/protoc",
	}
	protocArgs := []string{
		"--go_out=plugins=grpc:/go/src/",
		"-I",
		"/go/src/",
	}
	if protoc {
		for _, path := range paths {
			args := append(protocArgs, path)
			out, err := exec.Command("protoc", args...).CombinedOutput()
			if err != nil {
				fmt.Println(args)
				fmt.Println(string(out))
				panic(err)
			}
		}
	} else {
		for _, path := range paths {
			args := append(append(dockerArgs, protocArgs...), path)
			out, err := exec.Command("docker", args...).CombinedOutput()
			if err != nil {
				fmt.Println(args)
				fmt.Println(string(out))
				panic(err)
			}
		}
	}
}
