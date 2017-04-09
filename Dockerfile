FROM golang:1.8

COPY . /go/src/github.com/vrgakos/livemigrate

RUN mkdir /app
RUN mkdir -p $GOPATH/src/github.com/docker \
    && cd $GOPATH/src/github.com/docker \
    && git clone http://github.com/vrgakos/docker
RUN cd /go/src/github.com/vrgakos/livemigrate/migrate-docker \
    && go get -v
RUN go install github.com/vrgakos/livemigrate/migrate-docker

# DO SOME CLEANUP
RUN rm -rf /go/src

WORKDIR /app

ENTRYPOINT /go/bin/migrate-docker