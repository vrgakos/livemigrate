FROM golang:1.8

COPY . /go/src/github.com/vrgakos/livemigrate/webapp

RUN go install github.com/vrgakos/livemigrate/webapp

ENTRYPOINT /go/bin/webapp