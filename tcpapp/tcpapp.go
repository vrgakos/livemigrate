package tcpapp

import (
	"time"
)

type MsgRequest struct {
	Time		time.Time
	Data		[]byte
}

type MsgResponse struct {
	Time		time.Time
	Data		[]byte
}