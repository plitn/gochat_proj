package main

import (
	"fmt"
	"io"
	"net/http"
	"golang.org/x/net/websocket"
	"time"
)

type Server struct {
	connections map[*websocket.Conn]bool
}

func CreateServer() *Server {
	return &Server{
		connections: make(map[*websocket.Conn]bool),
	}
}

// TODO: add mutex
func (s *Server) WSHandler(ws *websocket.Conn) {
	fmt.Println("new connection from client at", ws.RemoteAddr())
	s.connections[ws] = true
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := ws.Read(buffer)
		if err != nil  && err != io.EOF {
			fmt.Println("read error", err)
			continue
		}
		msg := buffer[:n]
		s.broadcast(msg)
	}
}

func (s *Server) handleFeed(ws *websocket.Conn) {
	fmt.Println("conn from feed", ws.RemoteAddr())
	for {
		data := "time now: " + string(time.Now().Format("15:04:05"))
		ws.Write([]byte(data))
		time.Sleep(time.Second * 2)
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.connections {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error", err)
			}
		}(ws)
	}
}

func main() {
	server := CreateServer()
	http.Handle("/ws", websocket.Handler(server.WSHandler))
	http.Handle("/feed", websocket.Handler(server.handleFeed))
	http.ListenAndServe(":3000", nil)
}