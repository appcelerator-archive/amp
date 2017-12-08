package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/mitchellh/go-homedir"
)

var (
	ug              = "0:0"
	repo            = "github.com/appcelerator/amp"
	dockerCmd       = "gosu root docker"
	toolsImage      = "appcelerator/amptools:1.13"
	localToolsImage = "amptools"
	dockerArgs      []string
)

func init() {
	homedir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = os.TempDir()
		if err := os.Mkdir(goPath+"/pkg", os.ModePerm); err != nil {
			if err := os.Chmod(goPath+"/pkg", os.ModePerm); err != nil {
				panic(err)
			}
		}
		fmt.Printf("No GOPATH. Using %s as temporary GOPATH.\n", goPath)
	} else {
		fmt.Println("Using existing GOPATH:", goPath)
	}

	if runtime.GOOS == "linux" {
		ug = fmt.Sprintf("%s:%s", strconv.Itoa(os.Getuid()), strconv.Itoa(os.Getgid()))
	}

	dockerArgs = []string{
		"run", "-t", "--rm", "--name", "amptools",
		"-u", ug,
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-v", fmt.Sprintf("%s/.ssh:/root/.ssh:ro", homedir),
		"-v", fmt.Sprintf("%s/.config/amp:/root/.config/amp:ro", homedir),
		"-v", fmt.Sprintf("%s:/go/src/%s", wd, repo),
		"-v", fmt.Sprintf("%s/pkg:/go/pkg", goPath),
		"-w", fmt.Sprintf("/go/src/%s", repo),
		"-e", fmt.Sprintf("DOCKER_CMD=%s", dockerCmd),
	}
	if runtime.GOOS == "linux" {
		dockerArgs = append(dockerArgs, []string{localToolsImage}...)
	} else {
		dockerArgs = append(dockerArgs, []string{toolsImage}...)
	}
}

// build a local image to avoid leaving files with broken permissions
func buildLocalToolsImage() {
	// build the local image "amptools" for the current user
	content := []byte(fmt.Sprintf("FROM %s\nRUN sed -i \"s/sudoer:x:[0-9]*:[0-9]*/sudoer:x:%s/\" /etc/passwd", toolsImage, ug))
	tmpdir, err := ioutil.TempDir("", "dockerbuild")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(tmpdir) // clean up

	dockerfile := filepath.Join(tmpdir, "Dockerfile")
	if err := ioutil.WriteFile(dockerfile, content, 0666); err != nil {
		panic(err)
	}

	// docker build -t amptools tmpdir
	cmd := "docker"
	args := []string{
		"build",
		"-t",
		"amptools",
		tmpdir,
	}

	runcmd(cmd, args)
}

func main() {
	args := []string{
		"make",
	}

	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}

	if runtime.GOOS == "linux" {
		buildLocalToolsImage()
	}

	cmd := "docker"
	args = append(dockerArgs, args...)

	runcmd(cmd, args)
}

func runcmd(cmd string, args []string) {
	proc := exec.Command(cmd, args...)

	stdout, err := proc.StdoutPipe()
	if err != nil {
		panic(err)
	}
	outscanner := bufio.NewScanner(stdout)
	go func() {
		for outscanner.Scan() {
			fmt.Printf("%s\n", outscanner.Text())
		}
	}()

	stderr, err := proc.StderrPipe()
	if err != nil {
		panic(err)
	}
	errscanner := bufio.NewScanner(stderr)
	go func() {
		for errscanner.Scan() {
			fmt.Fprintf(os.Stderr, "%s\n", errscanner.Text())
		}
	}()

	err = proc.Start()
	if err != nil {
		panic(err)
	}

	err = proc.Wait()
	if err != nil {
		// Just pass along the information that the process exited with a failure;
		// whatever error information it displayed is what the user will see.
		// TODO: return the process exit code
		os.Exit(1)

	}
}
