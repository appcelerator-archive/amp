// Copyright 2015 David Chen <chendahui007@gmail.com>.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package safe implements a password safety checker.
package safe

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Level identifies the strength of the password.
type Level uint8

const (
	Terrible Level = iota
	Simple
	Medium
	Strong
)

// qwertyuiop
//  asdfghjkl
//   zxcvbnm
var (
	asdf    = "qwertyuiopasdfghjklzxcvbnm"
	revAsdf = reverse(asdf)

	words = make(map[string]int) // most frequent used passwords
)

var (
	lower  = regexp.MustCompile(`[a-z]`)
	upper  = regexp.MustCompile(`[A-Z]`)
	number = regexp.MustCompile(`[0-9]`)
	marks  = regexp.MustCompile(`[^0-9a-zA-Z]`)
)

type Safety struct {
	ml    int   // minimal length
	mf    int   // minimal frequency
	mt    int   // minimal type (minimum character family)
	level Level // default level to validate password
}

// New returns a Safety object.
func New(ml, mf, mt int, level Level) *Safety {
	return &Safety{ml, mf, mt, level}
}

func (s *Safety) SetWords(fpath string) error {
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	var tmp []string
	for scanner.Scan() {
		tmp = strings.Split(scanner.Text(), " ")
		words[tmp[0]], _ = strconv.Atoi(tmp[1])
	}
	return nil
}

// Check validate the password, and get the level.
func (s *Safety) Check(raw string) Level {
	l := len([]rune(raw))
	if l < s.ml {
		return Terrible
	}

	if s.isAsdf(raw) || s.isByStep(raw) {
		return Simple
	}

	if s.isCommonPassword(raw, s.mf) {
		return Simple
	}

	typ := 0

	if lower.MatchString(raw) {
		typ++
	}
	if upper.MatchString(raw) {
		typ++
	}
	if number.MatchString(raw) {
		typ++
	}
	if marks.MatchString(raw) {
		typ++
	}

	if l < 8 && typ == 2 {
		return Simple
	}

	if typ < s.mt {
		return Medium
	}

	return Strong
}

// If the password is in the order on keyboard.
func (s *Safety) isAsdf(raw string) bool {
	// s in asdf , or reverse in asdf
	rev := reverse(raw)
	return strings.Contains(asdf, raw) || strings.Contains(asdf, rev)
}

// If the password is alphabet step by step.
func (s *Safety) isByStep(raw string) bool {
	r := []rune(raw)
	delta := r[1] - r[0]

	for i, _ := range r {
		if i == 0 {
			continue
		}
		if r[i]-r[i-1] != delta {
			return false
		}
	}

	return true
}

// If the password is common used
// 10k top passwords: https://xato.net/passwords/more-top-worst-passwords/
func (s *Safety) isCommonPassword(raw string, freq int) bool {
	f, ok := words[raw]
	if !ok {
		return false
	}
	return f > freq
}

func reverse(raw string) string {
	r := []rune(raw)

	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
