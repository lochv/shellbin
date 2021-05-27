package agent

type Agent interface {
	read()
	write()
	Token() string
	GetReadChan() chan []byte
	GetWriteChan() chan []byte
	Disconnect()
}

func Listen(tcpPort string, udpPort string, aChan chan Agent) {
	tcpL := tcpListener{
		port:  tcpPort,
		aChan: aChan,
	}
	udpL := udpListener{
		port:  udpPort,
		aChan: aChan,
	}
	go tcpL.Listen()
	go udpL.Listen()
}
