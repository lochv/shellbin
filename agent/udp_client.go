package agent

import (
	"net"
)

type udpClient struct {
	conn      *net.UDPAddr
	token     string
	readChan  chan []byte
	writeChan chan []byte
}

func (u udpClient) Token() string {
	return u.token
}

func (u udpClient) GetReadChan() chan []byte {
	return u.readChan
}

func (u udpClient) GetWriteChan() chan []byte {
	return u.writeChan
}

func (u udpClient) Disconnect() {
	//TODO
}

func (u udpClient) read() {
	//TODO
}

func (u udpClient) write() {
	//TODO
}
