package agent

import (
	"net"
)

type tcpClient struct {
	conn      net.Conn
	token     string
	readChan  chan []byte
	writeChan chan []byte
}

func (t tcpClient) Token() string {
	return t.token
}

func (t tcpClient) GetReadChan() chan []byte {
	return t.readChan
}

func (t tcpClient) GetWriteChan() chan []byte {
	return t.writeChan
}

func (t tcpClient) Disconnect() {
	t.conn.Close()
}

func (t tcpClient) read() {
	for {
		buf := make([]byte, 1024)
		length, err := t.conn.Read(buf)
		if err != nil {
			t.conn.Close()
			return
		}
		if length > 0 {
			t.readChan <- buf[:length]
		}
	}
}

func (t tcpClient) write() {
	for {
		msg, ok := <-t.writeChan
		if !ok {
			t.conn.Close()
			return
		}
		t.conn.Write(msg)
	}
}
