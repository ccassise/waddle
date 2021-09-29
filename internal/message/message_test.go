package message

import (
	"testing"
)

func TestEqual(t *testing.T) {
	t.Run("should be equal", func(t *testing.T) {
		lhs := Message{
			Command:  Login,
			Receiver: "alice",
			Data:     "hello",
		}

		rhs := Message{
			Command:  Login,
			Receiver: "alice",
			Data:     "hello",
		}

		if !lhs.Equal(&rhs) {
			t.Fatalf("%v != %v", lhs, rhs)
		}
	})

	t.Run("should not be equal", func(t *testing.T) {
		lhs := Message{
			Command:  Login,
			Receiver: "alice",
			Data:     "hello",
		}

		rhs := Message{
			Command:  Join,
			Receiver: "bob",
			Data:     "goodbye",
		}

		if lhs.Equal(&rhs) {
			t.Fatalf("%v != %v", lhs, rhs)
		}
	})
}
