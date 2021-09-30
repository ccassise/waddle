package context

import (
	"testing"

	"github.com/ccassise/waddle/internal/message"
	"github.com/ccassise/waddle/internal/wuser"
	"github.com/ccassise/waddle/test/mock"
)

func TestLogin(t *testing.T) {
	t.Run("should login", func(t *testing.T) {
		ctx := New()
		user := wuser.User{Id: "alice_unique"}

		err := ctx.Login(&user, &message.Message{Data: "alice"})

		if err != nil || !user.LoggedIn || user.Name != "alice" {
			t.Fatalf("Login(%v) = (%v %v), want (%v, true)", user, err, user.LoggedIn, nil)
		}
	})

	t.Run("should fail to login on subsequent attempts", func(t *testing.T) {
		ctx := New()
		user := wuser.User{Id: "alice_unique"}

		ctx.Login(&user, &message.Message{Data: "alice"})
		err := ctx.Login(&user, &message.Message{Data: "new_alice"})

		if err == nil || !user.LoggedIn || user.Name != "alice" {
			t.Fatalf("Login(%v) = (%v %v), want (error, true)", user, err, user.LoggedIn)
		}
	})

	t.Run("should fail when name is already in use", func(t *testing.T) {
		ctx := New()
		users := []wuser.User{
			{
				Id: "alice_unique",
			},
			{
				Id: "bob_unique",
			},
		}

		ctx.Login(&users[0], &message.Message{Data: "alice"})
		err := ctx.Login(&users[1], &message.Message{Data: "alice"})

		if err == nil {
			t.Fatalf("Login(%v) = %q, want error", users[1], err)
		}
	})
}

func TestJoin(t *testing.T) {
	t.Run("should fail when not logged in", func(t *testing.T) {
		ctx := New()
		u := wuser.User{Id: "alice_unique"}

		err := ctx.Join(&u, &message.Message{Data: "#room"})

		if err == nil {
			t.Fatalf("Join() = %v, want %v;", err, nil)
		}
	})
}

func TestBroadcast(t *testing.T) {
	t.Run("should receive message when sent to chatroom", func(t *testing.T) {
		ctx := New()
		m := mock.MockWriter{Wrote: make([]byte, 0)}
		u := wuser.User{Id: "alice_unique", Writer: &m}

		ctx.Login(&u, &message.Message{Command: message.Login, Data: "alice"})
		ctx.Join(&u, &message.Message{Data: "#room"})
		err := ctx.Broadcast(&u, &message.Message{Command: message.Msg, Receiver: "#room", Data: "hello, room!"})

		expect := "GOTROOMMSG alice #room hello, room!\r\n"
		if err != nil || string(m.Wrote) != expect {
			t.Fatalf("Broadcast() = %v, want %v; sent %#q, want %#q", err, nil, string(m.Wrote), expect)
		}
	})

	t.Run("should receive message when sent to user", func(t *testing.T) {
		ctx := New()
		aliceWriter := mock.MockWriter{Wrote: make([]byte, 0)}
		bobWriter := mock.MockWriter{Wrote: make([]byte, 0)}
		alice := wuser.User{Id: "alice_unique", Writer: &aliceWriter}
		bob := wuser.User{Id: "bob_unique", Writer: &bobWriter}

		ctx.Login(&alice, &message.Message{Command: message.Login, Data: "alice"})
		ctx.Login(&bob, &message.Message{Command: message.Login, Data: "bob"})

		err := ctx.Broadcast(&alice, &message.Message{Command: message.Msg, Receiver: "bob", Data: "hello, bob!"})

		expect := "GOTUSERMSG alice hello, bob!\r\n"
		if err != nil || string(bobWriter.Wrote) != expect {
			t.Fatalf("Broadcast() = %v, want %v; sent %#q, want %#q", err, nil, string(bobWriter.Wrote), expect)
		}
	})

	t.Run("should fail when not logged in", func(t *testing.T) {
		ctx := New()
		m := mock.MockWriter{Wrote: make([]byte, 0)}
		u := wuser.User{Id: "alice_unique", Writer: &m}

		err := ctx.Broadcast(&u, &message.Message{Command: message.Msg, Receiver: "#room", Data: "hello, room!"})

		if err == nil {
			t.Fatalf("Broadcast() = %v, want %v; sent %#q, want %#q", err, nil, string(m.Wrote), "")
		}
	})

	t.Run("should receive message when from other user and only in that chatroom", func(t *testing.T) {
		ctx := New()
		aliceWriter := mock.MockWriter{Wrote: make([]byte, 0)}
		bobWriter := mock.MockWriter{Wrote: make([]byte, 0)}
		alice := wuser.User{Id: "alice_unique", Writer: &aliceWriter}
		bob := wuser.User{Id: "bob_unique", Writer: &bobWriter}

		ctx.Login(&alice, &message.Message{Command: message.Login, Data: "alice"})
		ctx.Login(&bob, &message.Message{Command: message.Login, Data: "bob"})
		ctx.Join(&alice, &message.Message{Data: "#room"})
		ctx.Join(&bob, &message.Message{Data: "#room"})
		ctx.Join(&bob, &message.Message{Data: "#test"})

		err := ctx.Broadcast(&alice, &message.Message{Command: message.Msg, Receiver: "#room", Data: "hello, room!"})

		expect := "GOTROOMMSG alice #room hello, room!\r\n"
		if err != nil || string(aliceWriter.Wrote) != expect || string(bobWriter.Wrote) != expect {
			t.Fatalf("Broadcast() = %v, want %v; got %#q, want %#q; got %#q, want %#q;",
				err, nil, string(bobWriter.Wrote), expect, string(bobWriter.Wrote), expect)
		}
	})

	t.Run("should fail when chatroom does not exist", func(t *testing.T) {
		ctx := New()

		m := mock.MockWriter{Wrote: make([]byte, 0)}
		u := wuser.User{Id: "alice_unique", Writer: &m}

		ctx.Login(&u, &message.Message{Command: message.Login, Data: "alice"})
		err := ctx.Broadcast(&u, &message.Message{Command: message.Msg, Receiver: "#room", Data: "hello, room!"})

		if err == nil {
			t.Fatalf("Broadcast() = %v, want error", err)
		}
	})

	t.Run("should fail when user not in chatroom", func(t *testing.T) {
		ctx := New()
		aliceWriter := mock.MockWriter{Wrote: make([]byte, 0)}
		bobWriter := mock.MockWriter{Wrote: make([]byte, 0)}
		alice := wuser.User{Id: "alice_unique", Writer: &aliceWriter}
		bob := wuser.User{Id: "bob_unique", Writer: &bobWriter}

		ctx.Login(&alice, &message.Message{Command: message.Login, Data: "alice"})
		ctx.Login(&bob, &message.Message{Command: message.Login, Data: "bob"})
		ctx.Join(&alice, &message.Message{Data: "#room"})
		ctx.Join(&bob, &message.Message{Data: "#test"})

		err := ctx.Broadcast(&alice, &message.Message{Command: message.Msg, Receiver: "#test", Data: "hello, room!"})

		if err == nil || string(bobWriter.Wrote) != "" {
			t.Fatalf("Broadcast() = %v, want error; sent %#q, want %#q", err, string(bobWriter.Wrote), "")
		}
	})

	t.Run("should fail when sending message to user not logged in", func(t *testing.T) {
		ctx := New()
		alice := wuser.User{Id: "alice_unique"}

		ctx.Login(&alice, &message.Message{Command: message.Login, Data: "alice"})

		err := ctx.Broadcast(&alice, &message.Message{Command: message.Msg, Receiver: "bob", Data: "hello, bob!"})

		if err == nil {
			t.Fatalf("Broadcast() = %v, want error", err)
		}
	})
}
