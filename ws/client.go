package ws

import (
	// "bytes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"shellbin/internal/logger"

	// "strings"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pingPeriod     = 30 * time.Second
	maxMessageSize = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  3024,
	WriteBufferSize: 3024,
}

type Client struct {
	conn      *websocket.Conn
	Token     string
	ReadChan  chan []byte
	WriteChan chan []byte
}

func (c Client) read(closeChan chan string) {
	defer func() {
		closeChan <- c.Token
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Write(err.Error())
			}
			break
		}
		c.ReadChan <- message
	}
}

func (c Client) write(closeChan chan string) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		closeChan <- c.Token
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.WriteChan:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.BinaryMessage, message)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func WsHandler(wsChan chan Client, closeChan chan string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	_, buff, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return
	}
	if len(string(buff)) != 16 {
		conn.Close()
		return
	}
	c := Client{
		conn:      conn,
		ReadChan:  make(chan []byte, 16),
		WriteChan: make(chan []byte, 16),
		Token:     string(buff),
	}
	wsChan <- c
	go c.read(closeChan)
	go c.write(closeChan)
}
