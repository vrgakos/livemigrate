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
	Time			time.Time
	Rtt			time.Duration
}

type TcpClient struct {
	target			string
	cb			TcpClientCb
	interval		time.Duration

	conn			net.Conn
	running			bool
	stopChan		chan bool
	stoppedChan		chan bool

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
	log.Println("TCP client: SendLoop started.")
	ticker := time.NewTicker(c.interval)
	//log.Printf("Task (%d) loop started\n", t.Id)

	dataLen := 1024
	data := make([]byte, dataLen)
	for i := 0; i < dataLen; i++ {
		data[i] = byte(i % 255)
	}

	for c.running {
		select {
		case <- c.stopChan:
			// Loop stopped
			break

		case <- ticker.C:
			// Tick
			req := &MsgRequest{
				Time:	time.Now(),
				Data:   data,
			}
			err := c.encoder.Encode(req)
			if err != nil {
				break
			}
		}
	}
	log.Println("TCP client: SendLoop ended.")
	c.stoppedChan <- true
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
			Time:	res.Time,
			Rtt:	time.Now().Sub(res.Time),
		}
		//log.Println("--RESULT", result)

		if c.cb != nil {
			c.cb(result)
		}
	}
	log.Println("TCP client: ReceiveLoop ended.")
	c.stoppedChan <- true
	c.Stop()
}

func (c *TcpClient) Start() error {
	if c.running {
		return fmt.Errorf("Already running!")
	}
	c.running = true
	c.stopChan = make(chan bool)
	c.stoppedChan = make(chan bool, 2)

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
	if !c.running {
		return fmt.Errorf("Not running!")
	}
	log.Println("TCP client: Stopping")
	//c.cb = nil
	c.running = false
	c.stopChan <- true

	// Close connection
	err := c.conn.Close()
	if err != nil {
		return err
	}

	<- c.stoppedChan
	<- c.stoppedChan
	log.Println("TCP client: Stopped")

	return nil
}
