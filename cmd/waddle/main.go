package main

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/ccassise/waddle/internal/context"
	"github.com/ccassise/waddle/internal/message"
	"github.com/ccassise/waddle/internal/parser"
	"github.com/ccassise/waddle/internal/wuser"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("usage: [port]")
	}

	port := ":" + os.Args[1]
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err.Error())
	}

	ctx := context.New()

	log.Println("Listening on", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		log.Printf("%v connect", conn.RemoteAddr())
		go handleConnection(&ctx, conn)
	}
}

func handleConnection(ctx *context.Context, conn net.Conn) {
	defer conn.Close()

	user := wuser.User{
		Id:     conn.RemoteAddr().String(),
		Writer: conn,
	}

	conn.Write([]byte("HELLO\r\n"))

	const maxMsgLen = 1024
	bytesRecv := make([]byte, maxMsgLen)
	for {
		n, err := conn.Read(bytesRecv)
		if err != nil {
			log.Printf("%v[%q] disconnect\n", user.Id, user.Name)
			return
		}

		msg, err := parser.Parse(string(bytesRecv[:n]))
		if err != nil {
			log.Printf("%v[%q] ERROR %q\n", user.Id, user.Name, err.Error())
			user.Error(err.Error())
			continue
		}

		log.Printf("%v[%q] %v %q %q\n", user.Id, user.Name, message.StringifyCommand(msg.Command), msg.Receiver, msg.Data)
		if err = execute(ctx, &user, &msg); err != nil {
			user.Error(err.Error())
			continue
		}

		user.Ok()
	}
}

func execute(ctx *context.Context, user *wuser.User, msg *message.Message) error {
	switch msg.Command {
	case message.Login:
		return executeLogin(ctx, user, msg)
	case message.Join:
		if !user.LoggedIn {
			return errors.New("unauthorized")
		}
		ctx.Join(msg.Data, user)
		return nil
	}
	return errors.New("internal error")
}

func executeLogin(ctx *context.Context, user *wuser.User, msg *message.Message) error {
	if err := ctx.Login(user); err != nil {
		log.Printf("%v[%q] login fail\n", user.Id, user.Name)
		return err
	}

	user.Name = msg.Data
	log.Printf("%v[%q] login successful\n", user.Id, user.Name)

	return nil
}
