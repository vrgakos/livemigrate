package tcpapp

import "time"

type MsgRequest struct {
	Time		time.Time
}

type MsgResponse struct {
	Time		time.Time
}