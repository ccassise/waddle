package main

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/ccassise/waddle/internal/context"
	"github.com/ccassise/waddle/internal/message"
	"github.com/ccassise/waddle/internal/parser"
	"github.com/ccassise/waddle/internal/wdluser"
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

	user := wdluser.User{
		Id:     conn.RemoteAddr().String(),
		Writer: conn,
	}
	defer ctx.Logout(&user)

	conn.Write([]byte("HELLO\r\n"))

	const maxBytes = 1024
	bytesRecv := make([]byte, maxBytes)
	for {
		n, err := conn.Read(bytesRecv)
		if err != nil {
			log.Printf("%v[%q] disconnect\n", user.Id, user.Name)
			return
		}

		msg, err := parser.Parse(bytesRecv[:n])
		if err != nil {
			log.Printf("%v[%q] ERROR %q\n", user.Id, user.Name, err.Error())
			user.Error(err.Error())
			continue
		}

		log.Printf("%v[%q] %v %q %q\n", user.Id, user.Name, message.StringifyCommand(msg.Command), msg.Receiver, msg.Data)
		if err = execute(ctx, &user, &msg); err != nil {
			log.Printf("%v[%q] ERROR %q\n", user.Id, user.Name, err.Error())
			user.Error(err.Error())
			continue
		}

		user.Ok()

		if msg.Command == message.Logout {
			break
		}
	}
}

func execute(ctx *context.Context, u *wdluser.User, m *message.Message) error {
	switch m.Command {
	case message.Login:
		return ctx.Login(u, m)
	case message.Logout:
		return ctx.Logout(u)
	case message.Join:
		return ctx.Join(u, m)
	case message.Part:
		return ctx.Part(u, m)
	case message.Msg:
		return ctx.Broadcast(u, m)
	}
	return errors.New("internal error")
}
