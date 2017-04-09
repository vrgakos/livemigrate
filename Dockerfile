FROM golang:1.8

ADD . /go/src/github.com/vrgakos/livemigrate

RUN mkdir /app
RUN mkdir -p cd $GOPATH/src/github.com/docker && cd $GOPATH/src/github.com/docker && git clone http://github.com/vrgakos/docker
RUN go get -v
RUN go install github.com/vrgakos/livemigrate/migrate-docker

WORKDIR /app

ENTRYPOINT /go/bin/migrate-docker