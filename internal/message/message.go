package message

import (
	"fmt"
)

// Message represents all of the necessary information in order to carry out a
// command.
type Message struct {
	Command  int
	Receiver string
	Data     string
}

// TODO: Add HELP command?
// List of commands.
const (
	Login = iota + 1
	Join
	Part
	Msg
	Logout
)

// Compares two messages and determines their equality.
func (t *Message) Equal(rhs *Message) bool {
	return t.Command == rhs.Command && t.Receiver == rhs.Receiver && t.Data == rhs.Data
}

func (t Message) String() string {
	return fmt.Sprintf("{ %v %#q %#q }", StringifyCommand(t.Command), t.Receiver, t.Data)
}

// Returns the string version of a given command. Used for testing/debugging.
func StringifyCommand(command int) string {
	switch command {
	case Login:
		return "LOGIN"
	case Join:
		return "JOIN"
	case Part:
		return "PART"
	case Msg:
		return "MSG"
	case Logout:
		return "LOGOUT"
	default:
		return ""
	}
}
