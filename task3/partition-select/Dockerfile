FROM golang:alpine

WORKDIR /app

RUN apk add --no-cache build-base zeromq-dev zeromq git mercurial \
    && go get -u github.com/golang/protobuf/proto \
    && go get -u github.com/pebbe/zmq4 \
    && apk del git mercurial

ADD . /src

RUN cd /src; go build -o pselect; cp pselect /app

CMD ["./pselect"]
