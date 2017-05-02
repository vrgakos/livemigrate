package tcpapp

import (
	"time"
	"github.com/gogo/protobuf/test/data"
)

type MsgRequest struct {
	Time		time.Time
	Data		[]byte
}

type MsgResponse struct {
	Time		time.Time
	Data		[]byte
}