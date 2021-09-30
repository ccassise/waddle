package parser

import (
	"testing"

	"github.com/ccassise/waddle/internal/message"
)

func TestParse(t *testing.T) {
	t.Run("LOGIN", func(t *testing.T) {
		input := []byte("LOGIN alice\r\n")

		actual, err := Parse(input)
		expect := message.Message{
			Command:  message.Login,
			Receiver: "",
			Data:     "alice",
		}

		if !actual.Equal(&expect) {
			t.Fatalf("Parse(%#q) = (%v, %v), want (%v, %v)", input, actual, err, expect, nil)
		}
	})

	t.Run("JOIN", func(t *testing.T) {
		t.Run("should succeed", func(t *testing.T) {
			input := []byte("JOIN #chatroom\r\n")

			actual, err := Parse(input)
			expect := message.Message{
				Command:  message.Join,
				Receiver: "",
				Data:     "#chatroom",
			}

			if !actual.Equal(&expect) {
				t.Fatalf("Parse(%#q) = (%v, %v), want (%v, %v)", input, actual, err, expect, nil)
			}
		})

		t.Run("should fail when missing '#'", func(t *testing.T) {
			input := []byte("JOIN chatroom\r\n")

			actual, err := Parse(input)
			expect := message.Message{}

			if !actual.Equal(&expect) || err == nil {
				t.Fatalf("Parse(%#q) = (%v, %v), want (%v, error)", input, actual, err, expect)
			}
		})

		t.Run("shoudl fail when missing <chatroom>", func(t *testing.T) {
			input := []byte("JOIN #\n")

			actual, err := Parse(input)
			expect := message.Message{}

			if !actual.Equal(&expect) || err == nil {
				t.Fatalf("Parse(%#q) = (%v, %v), want (%v, error)", input, actual, err, expect)
			}
		})
	})

	t.Run("MSG", func(t *testing.T) {
		t.Run("should succeed when chatroom", func(t *testing.T) {
			input := []byte("MSG #chatroom hello, world\r\n")

			actual, err := Parse(input)
			expect := message.Message{
				Command:  message.Msg,
				Receiver: "#chatroom",
				Data:     "hello, world",
			}

			if !actual.Equal(&expect) {
				t.Fatalf("Parse(%#q) = (%v, %v), want (%v, %v)", input, actual, err, expect, nil)
			}
		})

		t.Run("should succeed when user", func(t *testing.T) {
			input := []byte("MSG username hello, world\r\n")

			actual, err := Parse(input)
			expect := message.Message{
				Command:  message.Msg,
				Receiver: "username",
				Data:     "hello, world",
			}

			if !actual.Equal(&expect) {
				t.Fatalf("Parse(%#q) = (%v, %v), want (%v, %v)", input, actual, err, expect, nil)
			}
		})
	})
}
