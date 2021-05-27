package main

import (
	"shellbin/agent"
	"shellbin/http"
	"shellbin/hub"
	"shellbin/internal/config"
	"shellbin/ws"
)

func main() {
	agentChan := make(chan agent.Agent, 16)
	wsClientChan := make(chan ws.Client, 16)
	closeChan := make(chan string, 16)
	agent.Listen(config.Conf.TcpPort, config.Conf.UdpPort, agentChan)
	go http.Listen(config.Conf.HttpPort, wsClientChan, closeChan)
	h := hub.NewHub(agentChan, wsClientChan, closeChan)
	h.Run()
}
