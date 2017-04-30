FROM golang:1.8

COPY . /go/src/github.com/vrgakos/livemigrate/tcpapp

RUN go install github.com/vrgakos/livemigrate/tcpapp/cmd

ENTRYPOINT /go/bin/cmd