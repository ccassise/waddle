package wuser

import (
	"fmt"
	"io"
)

type User struct {
	Id       string
	Name     string
	LoggedIn bool
	Writer   io.Writer
}

// Writes OK to user. Return writer error.
func (u *User) Ok() error {
	_, err := u.Writer.Write([]byte("OK\r\n"))
	return err
}

// Writes ERROR and reason to user. Returns writer error.
func (u *User) Error(s string) error {
	_, err := u.Writer.Write([]byte(fmt.Sprintf("ERROR %v\r\n", s)))
	return err
}

// user.Msg(msg Message)
