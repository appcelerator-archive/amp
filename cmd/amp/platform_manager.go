package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/appcelerator/amp/config"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
)

const (
	//DockerURL docker url
	DockerURL = amp.DockerDefaultURL
	//DockerVersion docker version
	DockerVersion = amp.DockerDefaultVersion
	//ClearScreen ANSI Escape code
	ClearScreen = "\033[2J\033[0;0H"
	//MoveCursorHome ANSI Escape code
	MoveCursorHome = "\033[0;0H"
)

type ampManager struct {
	docker      *client.Client
	verbose     bool
	quiet       bool
	printColor  [6]*color.Color
	fcolRegular func(...interface{}) string
	fcolInfo    func(...interface{}) string
	fcolWarn    func(...interface{}) string
	fcolError   func(...interface{}) string
	fcolSuccess func(...interface{}) string
	fcolUser    func(...interface{}) string
	fcolTitle   func(...interface{}) string
	fcolLines   func(...interface{}) string
}

var currentColorTheme = "default"
var (
	colRegular = 0
	colInfo    = 1
	colWarn    = 2
	colError   = 3
	colSuccess = 4
	colUser    = 5
)

func newManager(verbose string) *ampManager {
	s := &ampManager{}
	s.setColors()
	if verbose == "true" {
		s.verbose = true
	}
	return s
}

func (s *ampManager) connectDocker() error {
	defaultHeaders := map[string]string{"User-Agent": "amplifier"}
	cli, err := client.NewClient(DockerURL, DockerVersion, nil, defaultHeaders)
	if err != nil {
		return fmt.Errorf("impossible to connect to Docker on: %s\n%v", DockerURL, err)
	}
	s.docker = cli
	return nil
}

func (s *ampManager) printf(col int, format string, args ...interface{}) {
	if s.quiet {
		return
	}
	colorp := s.printColor[0]
	if col > 0 && col < len(s.printColor) {
		colorp = s.printColor[col]
	}
	if !s.verbose && col == colInfo {
		return
	}
	colorp.Printf(format, args...)
}

func (s *ampManager) fatalf(format string, args ...interface{}) {
	s.printf(colError, format, args...)
	os.Exit(1)
}

func (s *ampManager) close() {
	//s.docker.Close()
}

func (s *ampManager) setColors() {
	//theme := s.getTheme()
	theme := AMP.Configuration.CmdTheme
	if theme == "dark" {
		s.printColor[0] = color.New(color.FgHiWhite)
		s.printColor[1] = color.New(color.FgHiBlack)
		s.printColor[2] = color.New(color.FgYellow)
		s.printColor[3] = color.New(color.FgRed)
		s.printColor[4] = color.New(color.FgGreen)
		s.printColor[5] = color.New(color.FgHiGreen)
	} else {
		s.printColor[0] = color.New(color.FgMagenta)
		s.printColor[1] = color.New(color.FgHiBlack)
		s.printColor[2] = color.New(color.FgYellow)
		s.printColor[3] = color.New(color.FgRed)
		s.printColor[4] = color.New(color.FgGreen)
		s.printColor[5] = color.New(color.FgHiGreen)
	} //add theme as you want.
	s.fcolRegular = s.printColor[colRegular].SprintFunc()
	s.fcolInfo = s.printColor[colInfo].SprintFunc()
	s.fcolWarn = s.printColor[colWarn].SprintFunc()
	s.fcolError = s.printColor[colError].SprintFunc()
	s.fcolSuccess = s.printColor[colSuccess].SprintFunc()
	s.fcolUser = s.printColor[colUser].SprintFunc()
	s.fcolTitle = s.printColor[colRegular].SprintFunc()
	s.fcolLines = s.printColor[colSuccess].SprintFunc()
}

func (s *ampManager) followClearScreen(follow bool) {
	if follow {
		fmt.Println(ClearScreen)
	}
}

func (s *ampManager) followMoveCursorHome(follow bool) {
	if follow {
		fmt.Println(MoveCursorHome)
	}
}

func (s *ampManager) displayInOrder(title1 string, title2 string, lines []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	sort.Strings(lines)
	if title1 != "" {
		fmt.Fprintln(w, s.fcolTitle(title1))
	}
	if title2 != "" {
		fmt.Fprintln(w, s.fcolTitle(title2))
	}
	for _, line := range lines {
		fmt.Fprintf(w, "%s\n", s.fcolLines(line))
	}
	w.Flush()
}

// system prerequisites
func (s *ampManager) systemPrerequisites() error {
	sysctl := false
	// checks if GOOS is set
	goos := os.Getenv("GOOS")
	if goos == "linux" {
		sysctl = true
	} else if goos == "" {
		// check if sysctl exists on the system
		if _, err := os.Stat("/etc/sysctl.conf"); err == nil {
			sysctl = true
		}
	}
	if sysctl {
		var out bytes.Buffer
		var stderr bytes.Buffer
		mmcmin := 262144
		cmd := exec.Command("sysctl", "-n", "vm.max_map_count")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		mmc, err := strconv.Atoi(strings.TrimRight(out.String(), "\n"))
		if err != nil {
			return err
		}
		if mmc < mmcmin {
			// admin rights are needed
			u, err := user.Current()
			if err != nil {
				return err
			}
			uid, err := strconv.Atoi(u.Uid)
			if err != nil {
				return err
			}
			if uid != 0 {
				return fmt.Errorf("vm.max_map_count should be at least 262144, admin rights are needed to update it")
			}
			if s.verbose {
				s.printf(colRegular, "setting max virtual memory areas\n")
			}
			cmd = exec.Command("sysctl", "-w", "vm.max_map_count=262144")
			err = cmd.Run()
		} else if s.verbose {
			s.printf(colRegular, "max virtual memory areas is already at a safe value\n")
		}
	}
	return nil
}
