package hub

import (
	"shellbin/agent"
	"shellbin/internal/logger"
	"shellbin/ws"
)

type hub struct {
	agents       map[string]agent.Agent
	wsClients    map[string]ws.Client
	agentChan    chan agent.Agent
	wsClientChan chan ws.Client
	closeChan    chan string
}

func NewHub(aChan chan agent.Agent, wsClientChan chan ws.Client, closeChan chan string) *hub {
	return &hub{
		agents:       make(map[string]agent.Agent),
		wsClients:    make(map[string]ws.Client),
		agentChan:    aChan,
		wsClientChan: wsClientChan,
		closeChan:    closeChan,
	}
}

func (h *hub) Run() {
	for {
		select {
		case agent := <-h.agentChan:
			logger.Write("new agent with token ", agent.Token())
			if wsClient, ok := h.wsClients[agent.Token()]; ok {
				h.agents[agent.Token()] = agent
				go func() {
					writeToAgentChan := agent.GetWriteChan()
					for {
						select {
						case msg := <-wsClient.ReadChan:
							logger.Write("client send ", msg)
							writeToAgentChan <- msg
						}
					}
				}()
				go func() {
					readFromAgentChan := agent.GetReadChan()
					for {
						select {
						case msg := <-readFromAgentChan:
							logger.Write("agent read ", msg)
							wsClient.WriteChan <- msg
						}
					}
				}()
			} else {
				agent.Disconnect()
			}
		case wsClient := <-h.wsClientChan:
			logger.Write("new ws client with token ", wsClient.Token)
			h.wsClients[wsClient.Token] = wsClient
		case token := <-h.closeChan:
			logger.Write("close ws ", token)
			_, ok := h.wsClients[token]
			if ok {
				delete(h.wsClients, token)
			}
			agent, ok := h.agents[token]
			if ok {
				agent.Disconnect()
				delete(h.agents, token)
			}
		}
	}
}
