package main

import (
	"os"

	"github.com/fatih/color"
)

type cmdManager struct {
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

func NewCmdManager(verbose string) *cmdManager {
	s := &cmdManager{}
	s.setColors()
	if verbose == "true" {
		s.verbose = true
	}
	return s
}

func (s *cmdManager) printf(col int, format string, args ...interface{}) {
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

func (s *cmdManager) fatalf(format string, args ...interface{}) {
	s.printf(colError, format, args...)
	os.Exit(1)
}

func (s *cmdManager) setColors() {
	//theme := s.getTheme()
	//theme := AMP.Configuration.CmdTheme
	// if theme == "dark" {
	// 	s.printColor[0] = color.New(color.FgHiWhite)
	// 	s.printColor[1] = color.New(color.FgHiBlack)
	// 	s.printColor[2] = color.New(color.FgYellow)
	// 	s.printColor[3] = color.New(color.FgRed)
	// 	s.printColor[4] = color.New(color.FgGreen)
	// 	s.printColor[5] = color.New(color.FgHiGreen)
	// } else {
	s.printColor[0] = color.New(color.FgMagenta)
	s.printColor[1] = color.New(color.FgHiBlack)
	s.printColor[2] = color.New(color.FgYellow)
	s.printColor[3] = color.New(color.FgRed)
	s.printColor[4] = color.New(color.FgGreen)
	s.printColor[5] = color.New(color.FgHiGreen)
	//} //add theme as you want.
	s.fcolRegular = s.printColor[colRegular].SprintFunc()
	s.fcolInfo = s.printColor[colInfo].SprintFunc()
	s.fcolWarn = s.printColor[colWarn].SprintFunc()
	s.fcolError = s.printColor[colError].SprintFunc()
	s.fcolSuccess = s.printColor[colSuccess].SprintFunc()
	s.fcolUser = s.printColor[colUser].SprintFunc()
	s.fcolTitle = s.printColor[colRegular].SprintFunc()
	s.fcolLines = s.printColor[colSuccess].SprintFunc()
}
