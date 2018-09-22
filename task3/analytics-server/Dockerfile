FROM golang:alpine

WORKDIR /app

RUN apk add --no-cache bash build-base zeromq-dev zeromq git mercurial \
    && go get -u github.com/go-sql-driver/mysql \
    && go get -u github.com/gorilla/mux \
    && go get -u github.com/golang/protobuf/proto \
    && go get -u github.com/pebbe/zmq4 \
    && apk del git mercurial

ADD . /src

RUN cd /src; go build -o server; cp server /app

CMD ["./server"]
