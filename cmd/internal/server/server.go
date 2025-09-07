package server

import (
	"fmt"
	"io"
	"log"
	"main.go/cmd/internal/request"
	"main.go/cmd/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	active   atomic.Bool
	handler  Handler
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (e *HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, e.StatusCode)
	messageBytes := []byte(e.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

func Serve(port int, handler Handler) (*Server, error) {
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: tcpListener,
		handler:  handler,
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
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}
	w := response.GetResponseWriter(&conn)
	s.handler(w, req)
	return
}
func (s *Server) Close() error {
	s.active.Store(false)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
