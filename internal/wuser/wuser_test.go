package wuser

import "testing"

type MockWriter struct {
	Wrote []byte
}

func (w *MockWriter) Write(b []byte) (int, error) {
	w.Wrote = append(w.Wrote, b...)
	return len(b), nil
}

func TestOk(t *testing.T) {
	m := MockWriter{make([]byte, 0)}
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
	m := MockWriter{make([]byte, 0)}
	u := User{
		Writer: &m,
	}

	err := u.Error("Hello, World")

	expect := "ERROR Hello, World\r\n"
	if err != nil || string(m.Wrote) != expect {
		t.Fatalf("Ok() = %#q, want %#q", m.Wrote, expect)
	}
}
