package cli

import (
	"os"

	"fmt"

	"github.com/fatih/color"
)

// CmdManager manages command output
type CmdManager struct {
	verbose     bool
	quiet       bool
	printColor  [5]*color.Color
	fcolRegular func(...interface{}) string
	fcolInfo    func(...interface{}) string
	fcolWarn    func(...interface{}) string
	fcolError   func(...interface{}) string
	fcolSuccess func(...interface{}) string
	fcolTitle   func(...interface{}) string
	fcolLines   func(...interface{}) string
}

var (
	colRegular = 0
	colInfo    = 1
	colWarn    = 2
	colError   = 3
	colSuccess = 4
)

func NewCmdManager(verbose bool) *CmdManager {
	s := &CmdManager{}
	s.setColors()
	if verbose {
		s.verbose = true
	}
	return s
}

func (s *CmdManager) printf(col int, format string, args ...interface{}) {
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
	fmt.Println("")
}

func (s *CmdManager) setColors() {
	theme := Configuration{}.CmdTheme
	if theme == "dark" {
		s.printColor[0] = color.New(color.FgHiWhite)
		s.printColor[1] = color.New(color.FgHiBlack)
		s.printColor[2] = color.New(color.FgYellow)
		s.printColor[3] = color.New(color.FgRed)
		s.printColor[4] = color.New(color.FgGreen)
	} else {
		s.printColor[0] = color.New(color.FgBlue)
		s.printColor[1] = color.New(color.FgHiBlack)
		s.printColor[2] = color.New(color.FgYellow)
		s.printColor[3] = color.New(color.FgRed)
		s.printColor[4] = color.New(color.FgGreen)
	} //add theme as you want.
	s.fcolRegular = s.printColor[colRegular].SprintFunc()
	s.fcolInfo = s.printColor[colInfo].SprintFunc()
	s.fcolWarn = s.printColor[colWarn].SprintFunc()
	s.fcolError = s.printColor[colError].SprintFunc()
	s.fcolSuccess = s.printColor[colSuccess].SprintFunc()
	s.fcolTitle = s.printColor[colRegular].SprintFunc()
	s.fcolLines = s.printColor[colSuccess].SprintFunc()
}

func (s *CmdManager) Regular(format string, args ...interface{}) {
	s.printf(colRegular, format, args...)
}

func (s *CmdManager) Info(format string, args ...interface{}) {
	s.printf(colInfo, format, args...)
}

func (s *CmdManager) Warn(format string, args ...interface{}) {
	s.printf(colWarn, format, args...)
}

func (s *CmdManager) Fatal(format string, args ...interface{}) {
	s.printf(colError, format, args...)
	os.Exit(1)
}

func (s *CmdManager) Success(format string, args ...interface{}) {
	s.printf(colSuccess, format, args...)
}
