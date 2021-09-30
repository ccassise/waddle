package context

import (
	"bytes"
	"errors"
	"strings"
	"sync"

	"github.com/ccassise/waddle/internal/message"
	"github.com/ccassise/waddle/internal/wuser"
)

// Context is a structure for shared data.
type Context struct {
	mu       sync.Mutex
	chatroom map[string][]*wuser.User
	user     map[string]*wuser.User
}

func New() Context {
	return Context{
		chatroom: make(map[string][]*wuser.User),
		user:     make(map[string]*wuser.User),
	}
}

// Login will login a user.
func (ctx *Context) Login(u *wuser.User, m *message.Message) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if u.LoggedIn {
		return errors.New(errUserLoggedIn)
	}

	if _, ok := ctx.user[m.Data]; ok {
		return errors.New(errUsernameInUse)
	}

	u.Name = m.Data
	u.LoggedIn = true
	ctx.user[u.Name] = u

	return nil
}

// Join will insert given user into given chatroom.
func (ctx *Context) Join(u *wuser.User, m *message.Message) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if !u.LoggedIn {
		return errors.New(errUnautorized)
	}

	room := m.Data

	if _, ok := ctx.chatroom[room]; !ok {
		ctx.chatroom[room] = []*wuser.User{u}
	} else {
		ctx.chatroom[room] = append(ctx.chatroom[room], u)
	}

	return nil
}

// Broadcast sends the given message from the given user to appropriate users.
func (ctx *Context) Broadcast(u *wuser.User, m *message.Message) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if !u.LoggedIn {
		return errors.New(errUnautorized)
	}

	if strings.HasPrefix(m.Receiver, "#") {
		return ctx.broadcastRoom(u, m)
	}

	return ctx.broadcastUser(u, m)
}

// broadcastRoom sends a given message from a given user to all users in a given room.
func (ctx *Context) broadcastRoom(u *wuser.User, m *message.Message) error {
	users, ok := ctx.chatroom[m.Receiver]
	if !ok {
		return errors.New(errUserNotInRoom)
	}

	isInRoom := false
	for i := range users {
		if users[i].Id == u.Id {
			isInRoom = true
			break
		}
	}

	if !isInRoom {
		return errors.New(errUserNotInRoom)
	}

	var buf bytes.Buffer
	buf.WriteString("GOTROOMMSG ")
	buf.WriteString(u.Name)
	buf.WriteString(" ")
	buf.WriteString(m.Receiver)
	buf.WriteString(" ")
	buf.WriteString(m.Data)
	buf.WriteString("\r\n")

	for i := range users {
		users[i].Writer.Write(buf.Bytes())
	}

	return nil
}

// broadcastUser sends a given message from a given user to a specific user.
func (ctx *Context) broadcastUser(u *wuser.User, m *message.Message) error {
	to, ok := ctx.user[m.Receiver]
	if !ok {
		return errors.New(errUserNotLoggedIn)
	}

	var buf bytes.Buffer
	buf.WriteString("GOTUSERMSG ")
	buf.WriteString(u.Name)
	buf.WriteString(" ")
	buf.WriteString(m.Data)
	buf.WriteString("\r\n")

	_, err := to.Writer.Write(buf.Bytes())
	if err != nil {
		return errors.New(errSendFailed)
	}

	return nil
}

const (
	errSendFailed      = "failed to send message"
	errUnautorized     = "unauthorized"
	errUserLoggedIn    = "user already logged in"
	errUserNotInRoom   = "user not in room"
	errUserNotLoggedIn = "user not logged in"
	errUsernameInUse   = "username already in use"
)
