FROM golang:1.8

COPY . /go/src/github.com/vrgakos/livemigrate/tcpapp

RUN apt update && apt install stress
RUN go install github.com/vrgakos/livemigrate/tcpapp/cmd

ENTRYPOINT /go/bin/cmd