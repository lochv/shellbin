package agent

import (
	"net"
)

type udpListener struct {
	port  string
	aChan chan Agent
}

func (u udpListener) Listen() {
 //TODO
}

func (u udpListener) newConnHandler(conn net.Conn) {
	//TODO
}
