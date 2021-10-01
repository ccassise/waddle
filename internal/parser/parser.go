package parser

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode"

	"github.com/ccassise/waddle/internal/message"
)

// Parse will parse the given string and return its representation.
func Parse(b []byte) (message.Message, error) {
	p := parser{
		msg: message.Message{},
		buf: bytes.NewBuffer(b),
	}

	ch, err := p.buf.ReadByte()
	if err != nil {
		return message.Message{}, err
	}

	switch ch {
	case 'J':
		p.buf.UnreadByte()
		err = p.parseOneArgCmd(message.Join, p.parseRoom)
		if err != nil {
			return message.Message{}, err
		}
	case 'L':
		p.buf.UnreadByte()
		err = p.parseL()
		if err != nil {
			return message.Message{}, err
		}
	case 'M':
		p.buf.UnreadByte()
		err = p.parseMessage()
		if err != nil {
			return message.Message{}, err
		}
	case 'P':
		p.buf.UnreadByte()
		err = p.parseOneArgCmd(message.Part, p.parseRoom)
		if err != nil {
			return message.Message{}, err
		}
	}

	return p.msg, nil
}

type parser struct {
	msg message.Message
	buf *bytes.Buffer
}

const (
	errChatroom       = "chatrooms must begin with '#'"
	errInvalidArgs    = "invalid arguments"
	errInvalidCommand = "invalid command"
)

// parseKeyword checks that the next word in the buffer matches exactly the
// given keyword.
func (p *parser) parseKeyword(keyword string) error {
	err := p.matchNext(keyword)
	if err != nil {
		return err
	}

	ch, err := p.buf.ReadByte()
	if err != nil {
		return err
	}

	if !unicode.IsSpace(rune(ch)) {
		return errors.New(errInvalidCommand)
	}

	return nil
}

// matchNext returns whether or not the next bytes in the buffer match exactly
// the given string.
func (p *parser) matchNext(word string) error {
	for i := range word {
		ch, err := p.buf.ReadByte()
		if err != nil {
			return err
		}

		if ch != word[i] {
			return errors.New(errInvalidCommand)
		}
	}

	return nil
}

// parseSpace skips all space characters.
func (p *parser) parseSpace() error {
	ch, err := p.buf.ReadByte()
	if err != nil {
		return err
	}

	for unicode.IsSpace(rune(ch)) {
		ch, err = p.buf.ReadByte()
		if err != nil {
			return err
		}
	}

	err = p.buf.UnreadByte()
	if err != nil {
		return err
	}

	return nil
}

// parseWord collects all symbols/letters until a space or error is found.
func (p *parser) parseWord() (string, error) {
	var result strings.Builder

	ch, err := p.buf.ReadByte()
	if err != nil {
		return "", err
	} else if unicode.IsSpace(rune(ch)) {
		return "", errors.New(errInvalidArgs)
	}

	for !unicode.IsSpace(rune(ch)) {
		err = result.WriteByte(ch)
		if err != nil {
			return "", err
		}

		ch, err = p.buf.ReadByte()
		if err != nil {
			return "", err
		}
	}

	err = p.buf.UnreadByte()
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

// #<chatroom>
func (p *parser) parseRoom() (string, error) {
	var result strings.Builder

	ch, err := p.buf.ReadByte()
	if err != nil {
		return "", err
	} else if ch != '#' {
		p.buf.UnreadByte()
		return "", errors.New(errChatroom)
	}

	_, err = result.WriteString("#")
	if err != nil {
		return "", err
	}

	word, err := p.parseWord()
	if err != nil {
		return "", err
	}

	_, err = result.WriteString(word)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

// parseMsgText reads all bytes until a carriage return or line feed is found.
func (p *parser) parseMsgText() (string, error) {
	var result strings.Builder

	for {
		ch, err := p.buf.ReadByte()
		if err != nil {
			return "", err
		}

		if ch == '\r' || ch == '\n' {
			break
		}

		err = result.WriteByte(ch)
		if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}

// parseOneArgCmd parses a <command> <argument> . The second function argument
// is a function that determines what strategy to use to parse the argument.
// E.G. #<chatroom> is parsed different than <username> .
func (p *parser) parseOneArgCmd(cmd int, parseArg func() (string, error)) error {
	err := p.parseKeyword(message.StringifyCommand(cmd))
	if err != nil {
		return err
	}

	p.msg.Command = cmd

	err = p.parseSpace()
	if err != nil {
		return err
	}

	p.msg.Data, err = parseArg()
	if err != nil {
		return err
	}

	err = p.parseSpace()
	if err != io.EOF {
		return errors.New(errInvalidArgs)
	}

	return nil
}

func (p *parser) parseMessage() error {
	err := p.parseKeyword(message.StringifyCommand(message.Msg))
	if err != nil {
		return err
	}

	p.msg.Command = message.Msg

	err = p.parseSpace()
	if err != nil {
		return err
	}

	p.msg.Receiver, err = p.parseRoom()
	if err != nil {
		p.msg.Receiver, err = p.parseWord()
		if err != nil {
			return err
		}
	}

	err = p.parseSpace()
	if err != nil {
		return err
	}

	p.msg.Data, err = p.parseMsgText()
	if err != nil {
		return err
	}

	err = p.parseSpace()
	if err != io.EOF {
		return errors.New(errInvalidArgs)
	}

	return nil
}

// parseL will determine if input is LOGIN or LOGOUT and parse correctly.
func (p *parser) parseL() error {
	err := p.matchNext("LOG")
	if err != nil {
		return errors.New(errInvalidCommand)
	}

	ch, err := p.buf.ReadByte()
	if err != nil {
		return err
	}

	if ch == 'O' {
		p.buf.UnreadByte()
		err = p.matchNext("OUT")
		if err != nil {
			return err
		}

		// Check next byte is not EOF.
		_, err = p.buf.ReadByte()
		if err != nil {
			return err
		}
		p.buf.UnreadByte()

		p.msg.Command = message.Logout
	} else {
		p.buf.UnreadByte()
		err = p.matchNext("IN")
		if err != nil {
			return err
		}

		p.msg.Command = message.Login

		err = p.parseSpace()
		if err != nil {
			return err
		}

		p.msg.Data, err = p.parseWord()
		if err != nil {
			return err
		}
	}

	err = p.parseSpace()
	if err != io.EOF {
		return errors.New(errInvalidCommand)
	}

	return nil
}
