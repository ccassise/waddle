package context

import (
	"testing"

	"github.com/ccassise/waddle/internal/wuser"
)

type WritableMock struct {
	Read []byte
}

func TestLogin(t *testing.T) {
	t.Run("should login", func(t *testing.T) {
		ctx := New()
		user := wuser.User{
			Id:   "alice_unique",
			Name: "alice",
		}

		err := ctx.Login(&user)

		if err != nil || !user.LoggedIn {
			t.Fatalf("Login(%v) = (%v %v), want (nil, true)", user, err, user.LoggedIn)
		}
	})

	t.Run("should fail to login on subsequent attempts", func(t *testing.T) {
		ctx := New()
		user := wuser.User{
			Id:   "alice_unique",
			Name: "alice",
		}

		ctx.Login(&user)
		err := ctx.Login(&user)

		if err == nil || !user.LoggedIn {
			t.Fatalf("Login(%v) = (%v %v), want (error, true)", user, err, user.LoggedIn)
		}
	})

	t.Run("should fail when name is already in use", func(t *testing.T) {
		ctx := New()
		users := []wuser.User{
			{
				Id:   "alice_unique",
				Name: "alice",
			},
			{
				Id:   "bob_unique",
				Name: "alice",
			},
		}

		ctx.Login(&users[0])
		err := ctx.Login(&users[1])

		if err == nil {
			t.Fatalf("Login(%v) = %q, want error", users[1], err)
		}
	})
}
