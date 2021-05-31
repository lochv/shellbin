package agent

import (
	"net"
	"shellbin/internal/logger"
)

type tcpListener struct {
	port  string
	aChan chan Agent
}

func (t tcpListener) Listen() {
	l, err := net.Listen("tcp", "0.0.0.0:"+t.port)
	if err != nil {
		panic(err.Error())
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			logger.Write(err.Error())
		}
		go t.newConnHandler(c)
	}
}

func (t tcpListener) newConnHandler(conn net.Conn) {
	buff := make([]byte, 16)
	n, err := conn.Read(buff)
	if err != nil {
		logger.Write(err.Error())
		conn.Close()
	}
	if n != 16 {
		conn.Close()
		return
	}
	c := tcpClient{
		conn:      conn,
		token:     string(buff),
		readChan:  make(chan []byte, 16),
		writeChan: make(chan []byte, 16),
	}
	go c.read()
	go c.write()

	t.aChan <- c
}
