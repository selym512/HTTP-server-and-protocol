package server

import (
	"fmt"
	"log"
	"main.go/cmd/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	active   atomic.Bool
}

func Serve(port int) (*Server, error) {
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: tcpListener,
	}
	go s.listen()
	return s, nil
}
func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.active.Load() == false {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}
func (server *Server) handle(conn net.Conn) {
	defer conn.Close()
	response.WriteStatusLine(conn, 200)
	headers := response.GetDefaultHeaders(0)
	response.WriteHeaders(conn, headers)
	return
}
func (s *Server) Close() error {
	s.active.Store(false)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
