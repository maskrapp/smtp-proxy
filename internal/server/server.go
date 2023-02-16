package server

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/pires/go-proxyproto"
)

type Server struct {
	mailserver *net.TCPAddr
	ln         net.Listener
}

func New(mailserver string) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", mailserver)
	if err != nil {
		return nil, err
	}
	return &Server{mailserver: addr}, nil
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", "0.0.0.0:25")
	if err != nil {
		panic(err)
	}
	s.ln = ln
	s.listen()
}

func (s *Server) Shutdown() error {
	return s.ln.Close()
}

func (s *Server) listen() {
	for {
		c, err := s.ln.Accept()
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}
		connection := c.(*net.TCPConn)
		connection.SetDeadline(time.Now().Add(2 * time.Minute))
		defer connection.Close()
		go func() {
			clientConn, err := net.DialTCP("tcp", nil, s.mailserver)
			if err != nil {
				fmt.Println("error: ", err)
				return
			}
			defer clientConn.Close()
			header := proxyproto.HeaderProxyFromAddrs(1, connection.RemoteAddr(), clientConn.RemoteAddr())
			if _, err := header.WriteTo(clientConn); err != nil {
				fmt.Println("proxy protocol error: ", err)
				return
			}
			s.startProxy(clientConn, connection)
		}()
	}
}

func (s *Server) startProxy(client, server *net.TCPConn) {
	serverShutdown := make(chan struct{}, 1)
	clientShutdown := make(chan struct{}, 1)

	go func() {
		io.Copy(server, client)
		clientShutdown <- struct{}{}
	}()
	go func() {
		io.Copy(client, server)
		serverShutdown <- struct{}{}
	}()
	select {
	case <-clientShutdown:
		server.CloseRead()
	case <-serverShutdown:
		client.CloseRead()
	}
}
