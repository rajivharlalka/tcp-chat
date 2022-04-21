package chat

import "net"

type client struct {
	conn     net.Addr
	nick     string
	room     *room
	commangs chan<- commandID
}
