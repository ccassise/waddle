package context

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ccassise/waddle/internal/wuser"
)

// type Chatroom struct {
// 	rooms map[string][]*wuser.User
// }

// chatroom.Join(room string, u *User)

// chatroom.Part(room string, u *User)

// chatroom.Broadcast(room string, msg Message)

type Context struct {
	mu       sync.Mutex
	chatroom map[string][]*wuser.User
	// user      map[string]*wuser.User
	nameInUse map[string]bool
}

func New() Context {
	return Context{
		chatroom:  make(map[string][]*wuser.User),
		nameInUse: make(map[string]bool),
	}
}

func (ctx *Context) Login(u *wuser.User) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if u.LoggedIn {
		return errors.New("user already logged in")
	}

	if isInUse, ok := ctx.nameInUse[u.Name]; ok && isInUse {
		return errors.New("username already in use")
	}

	u.LoggedIn = true
	ctx.nameInUse[u.Name] = true

	return nil
}

func (ctx *Context) Join(room string, u *wuser.User) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if _, ok := ctx.chatroom[room]; !ok {
		ctx.chatroom[room] = make([]*wuser.User, 1)
	}

	ctx.chatroom[room] = append(ctx.chatroom[room], u)
	fmt.Println(ctx.chatroom)
}

// context.Logout(u *User)
