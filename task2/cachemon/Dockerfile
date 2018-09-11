FROM golang:alpine

ADD . /go/src/app
WORKDIR /go/src/app

RUN apk add --no-cache git mercurial \
    && go get -u github.com/streadway/amqp \
    && apk del git mercurial

CMD ["go", "run", "mon.go"]
