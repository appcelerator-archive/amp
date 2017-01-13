// Copyright 2015 David Chen <chendahui007@gmail.com>.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package safe

import "testing"

var s = New(8, 0, 3, Strong)

func setup(t *testing.T) {
	err := s.SetWords("./words.dat")
	if err != nil {
		t.Errorf("set word error: %s\n", err.Error())
	}
}

func TestCheck(t *testing.T) {
	setup(t)
	for _, c := range []struct {
		in   string
		want Level
	}{
		{"dasda", Terrible},
		{"aW1#", Terrible},
		{"asdfghjkl", Simple},
		{"abcdefghi", Simple},
		{"password", Simple},
		{"qeasdasddsad", Medium},
		{"eqweqwewe123", Medium},
		{"asdQWEaaa", Medium},
		{"ewqeqwewqe12#", Strong},
	} {
		got := s.Check(c.in)
		if got != c.want {
			t.Errorf("%s got %v, want %v", c.in, got, c.want)
		}
	}

}

func TestIsAsdf(t *testing.T) {
	for _, c := range []struct {
		in   string
		want bool
	}{
		{"qwer", true},
		{"tyuio", true},
		{"asdf", true},
		{"lkjhg", true},
		{"zxcvb", true},
		{"mnbvc", true},
		{"Asdf", false},
		{"qwrty", false},
		{"lkjhgfdsa", true},
	} {
		got := s.isAsdf(c.in)
		if got != c.want {
			t.Errorf("got %t want %t", got, c.want)
		}
	}
}

func TestIsByStep(t *testing.T) {
	for _, c := range []struct {
		in   string
		want bool
	}{
		{"abc", true},
		{"hijklmn", true},
		{"aceg", true},
		{"asdf", false},
		{"123456", true},
		{"13579", true},
		{"123567", false},
	} {
		got := s.isByStep(c.in)
		if got != c.want {
			t.Errorf("got %t want %t", got, c.want)
		}
	}

}

func TestIsCommonPassword(t *testing.T) {
	setup(t)

	for _, c := range []struct {
		in   string
		freq int
		want bool
	}{
		{"password", 0, true},
		{"boat", 200, true},
		{"engine", 216, false},
		{"golang", 0, false},
	} {
		got := s.isCommonPassword(c.in, c.freq)
		if got != c.want {
			t.Errorf("got %t want %t", got, c.want)
		}
	}
}

func BenchmarkIsAsdf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.isAsdf("asdfghjkl")
	}
}

func BenchmarkIsByStep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.isByStep("abcdefg")
	}
}

func BenchmarkIsCommonPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.isCommonPassword("password", 0)
	}
}

func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reverse("qwertyasdfghjklmnbvcxz")
	}
}

func BenchmarkCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.Check("qwert123!Z")
	}
}
