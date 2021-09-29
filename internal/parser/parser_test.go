package parser

import (
	"testing"

	"github.com/ccassise/waddle/internal/message"
)

func TestParseLOGIN(t *testing.T) {
	input := "LOGIN alice\r\n"

	actual, err := Parse(input)
	expect := message.Message{
		Command:  message.Login,
		Receiver: "",
		Data:     "alice",
	}

	if !actual.Equal(&expect) {
		t.Fatalf("Parse(%#q) = (%v, %v), want (%v, nil)", input, actual, err, expect)
	}
}

func TestParseJOIN(t *testing.T) {
	test1 := "JOIN #chatroom\r\n"
	t.Run(test1, func(t *testing.T) {
		actual, err := Parse(test1)
		expect := message.Message{
			Command:  message.Join,
			Receiver: "",
			Data:     "#chatroom",
		}

		if !actual.Equal(&expect) {
			t.Fatalf("Parse(%#q) = (%v, %v), want (%v, nil)", test1, actual, err, expect)
		}
	})

	test2 := "JOIN chatroom\r\n"
	t.Run(test2, func(t *testing.T) {
		actual, err := Parse(test2)
		expect := message.Message{}

		if !actual.Equal(&expect) || err == nil {
			t.Fatalf("Parse(%#q) = (%v, %v), want (%v, error)", test2, actual, err, expect)
		}
	})

	test3 := "JOIN #\n"
	t.Run(test3, func(t *testing.T) {
		actual, err := Parse(test3)
		expect := message.Message{}

		if !actual.Equal(&expect) || err == nil {
			t.Fatalf("Parse(%#q) = (%v, %v), want (%v, error)", test3, actual, err, expect)
		}
	})
}
