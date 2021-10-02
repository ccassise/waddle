package wdluser

import (
	"testing"

	"github.com/ccassise/waddle/test/mock"
)

func TestOk(t *testing.T) {
	m := mock.MockWriter{Wrote: make([]byte, 0)}
	u := User{
		Writer: &m,
	}

	err := u.Ok()

	expect := "OK\r\n"
	if err != nil || string(m.Wrote) != expect {
		t.Fatalf("Ok() = %#q, want %#q", m.Wrote, expect)
	}
}

func TestError(t *testing.T) {
	m := mock.MockWriter{Wrote: make([]byte, 0)}
	u := User{
		Writer: &m,
	}

	err := u.Error("Hello, World")

	expect := "ERROR Hello, World\r\n"
	if err != nil || string(m.Wrote) != expect {
		t.Fatalf("Ok() = %#q, want %#q", m.Wrote, expect)
	}
}
