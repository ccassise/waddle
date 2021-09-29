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
func Parse(line string) (message.Message, error) {
	p := parser{
		msg: message.Message{},
		buf: bytes.NewBufferString(line),
	}
	ch, err := p.buf.ReadByte()
	if err != nil {
		return message.Message{}, err
	}

	switch ch {
	case 'L':
		err = p.buf.UnreadByte()
		if err != nil {
			return message.Message{}, err
		}

		err = p.parseOneArgCmd(message.Login, p.parseWord)
		if err != nil {
			return message.Message{}, err
		}
	case 'J':
		err = p.buf.UnreadByte()
		if err != nil {
			return message.Message{}, err
		}

		err = p.parseOneArgCmd(message.Join, p.parseRoom)
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
	for i := range keyword {
		ch, err := p.buf.ReadByte()
		if err != nil {
			return err
		}

		if ch != keyword[i] {
			return errors.New(errInvalidCommand)
		}
	}

	ch, err := p.buf.ReadByte()
	if err != nil {
		return err
	}

	if !unicode.IsSpace(rune(ch)) {
		return errors.New(errInvalidCommand)
	}

	err = p.buf.UnreadByte()
	if err != nil {
		return err
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

// func (p *parser) parseCRLF() error {
// 	const CRLF = "\r\n"
// 	for i := range CRLF {
// 		ch, err := p.buf.ReadByte()
// 		if err != nil {
// 			return err
// 		}

// 		if ch != CRLF[i] {
// 			return errors.New(errInvalidCommand)
// 		}
// 	}

// 	return nil
// }

// parseOneArgCmd parses a <command> <argument> . The second function argument
// is a function that determines what strategy to use to parse the argument.
// E.G. #<chatroom> is parse different than <username> .
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

	word, err := parseArg()
	if err != nil {
		return err
	}

	p.msg.Data = word

	err = p.parseSpace()
	if err != io.EOF {
		return errors.New(errInvalidArgs)
	}

	return nil
}
