package tcpapp

import (
	"net"
	"encoding/gob"
	"log"
	"fmt"
)

type TcpServer struct {
	target			string
	ln			net.Listener

	running			bool
}

func NewTcpServer(target string) *TcpServer {
	server := &TcpServer{
		target:		target,

		running:	false,
	}

	//client.Start()
	return server
}

func (s *TcpServer) Loop() {
	for s.running {
		conn, err := s.ln.Accept()
		if err != nil {
			continue
		}

		go s.ReceiveLoop(conn)
	}
}

func (s *TcpServer) ReceiveLoop(conn net.Conn) {
	log.Printf("New client: %s\n", conn.RemoteAddr().String())
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	for s.running {
		req := &MsgRequest{}
		err := decoder.Decode(req)
		if err != nil {
			break
		}

		//log.Println("Got REQUEST", req)
		res := &MsgResponse{
			Time:	req.Time,
		}
		err = encoder.Encode(res)
		if err != nil {
			break
		}
	}
	log.Printf("Client left: %s\n", conn.RemoteAddr().String())
	conn.Close()
}

func (s *TcpServer) Start() error {
	if s.running {
		return fmt.Errorf("Already running!")
	}
	s.running = true

	ln, err := net.Listen("tcp", s.target)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.Loop()

	return nil
}

func (s *TcpServer) Stop() error {
	if s.running {
		return fmt.Errorf("Not running!")
	}
	s.running = false

	// Close connection
	err := s.ln.Close()
	if err != nil {
		return err
	}

	return nil
}
