package main

import (
	"fmt"
	"log"
	"net"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has connected: %s", conn.RemoteAddr().String)

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) nick(c *client, args []string) {
	c.nick = args[1]
	c.msg(fmt.Sprintf("All Right,I will call you, %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	roomName := args[1]
	r, ok := s.rooms[roomName]

	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}

		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)

	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))
	c.msg(fmt.Sprintf("Welcome to %s", r.name))
}

func (s *server) listRooms(c *client, args []string) {}

func (s *server) msg(c *client, args []string) {}

func (s *server) quit(c *client, args []string) {}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
