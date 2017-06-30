package ampadmin

import (
	"testing"
	"os"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
}

func teardown() {
}

func TestFoo(t *testing.T) {
	if Foo() != "foo" {
		t.Error("Foo() != \"foo\"")
	}
}
