package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"shellbin/ws"
)

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "template/index.html")
}

func terminalHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "template/terminal.html")
}

func Listen(port string, wsChan chan ws.Client, closeChan chan string) {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/terminal", terminalHandler)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { ws.WsHandler(wsChan, closeChan, w, r) })
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("template/"))))
	http.ListenAndServe(":"+port, router)
}
