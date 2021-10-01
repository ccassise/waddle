package wuser

import (
	"bytes"
	"io"
)

type User struct {
	Id       string
	Name     string
	LoggedIn bool
	Writer   io.Writer
	Rooms    []string
}

// Writes OK to user. Return writer error.
func (u *User) Ok() error {
	_, err := u.Writer.Write([]byte("OK\r\n"))
	return err
}

// Writes ERROR and reason to user. Returns writer error.
func (u *User) Error(s string) error {
	var buf bytes.Buffer
	buf.WriteString("ERROR ")
	buf.WriteString(s)
	buf.WriteString("\r\n")

	_, err := u.Writer.Write(buf.Bytes())
	return err
}
