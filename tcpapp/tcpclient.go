package tcpapp

import (
	"net"
	"fmt"
	"encoding/gob"
	"time"
	"log"
)

type TcpClientCb func(*TcpClientResult)

type TcpClientResult struct {
	Rtt			time.Duration
}

type TcpClient struct {
	target			string
	cb			TcpClientCb
	interval		time.Duration

	conn			net.Conn
	running			bool
	stopChan		chan bool

	encoder			*gob.Encoder
	decoder			*gob.Decoder
}

func NewTcpClient(target string, interval time.Duration, cb TcpClientCb) *TcpClient {
	client := &TcpClient{
		target:		target,
		cb:		cb,
		interval:	interval,

		running:	false,
	}

	//client.Start()
	return client
}

func (c *TcpClient) Loop() {
	log.Println("Loop started.")
	ticker := time.NewTicker(c.interval)
	//log.Printf("Task (%d) loop started\n", t.Id)
	for c.running {
		select {
		case <- c.stopChan:
			// Loop stopped
			break

		case <- ticker.C:
			// Tick
			req := &MsgRequest{
				Time:	time.Now(),
			}
			err := c.encoder.Encode(req)
			if err != nil {
				break
			}
		}
	}
	log.Println("Loop ended.")
	c.Stop()
}

func (c *TcpClient) ReceiveLoop() {
	for c.running {
		res := &MsgResponse{}
		err := c.decoder.Decode(res)
		if err != nil {
			log.Println(err)
			break
		}

		//log.Println("Got RESPONSE", res)
		result := &TcpClientResult{
			Rtt:	time.Now().Sub(res.Time),
		}
		//log.Println("--RESULT", result)

		if c.cb != nil {
			c.cb(result)
		}
	}
	log.Println("ReceiveLoop ended.")
	c.Stop()
}

func (c *TcpClient) Start() error {
	if c.running {
		return fmt.Errorf("Already running!")
	}
	c.running = true
	c.stopChan = make(chan bool)

	conn, err := net.Dial("tcp", c.target)
	if err != nil {
		return err
	}
	c.conn = conn

	c.encoder = gob.NewEncoder(conn)
	c.decoder = gob.NewDecoder(conn)

	go c.ReceiveLoop()
	go c.Loop()

	return nil
}

func (c *TcpClient) Stop() error {
	log.Println("Stopping")
	if c.running {
		return fmt.Errorf("Not running!")
	}
	c.running = false
	c.stopChan <- true

	// Close connection
	err := c.conn.Close()
	if err != nil {
		return err
	}

	return nil
}
