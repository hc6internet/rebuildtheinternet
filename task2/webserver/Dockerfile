FROM golang:alpine

ADD . /go/src/app
WORKDIR /go/src/app

RUN apk add --no-cache bash git mercurial \
    && go get -u github.com/bradfitz/gomemcache/memcache \
    && go get -u github.com/streadway/amqp \
    && apk del git mercurial

# add memcached
RUN adduser -D memcache \
    && apk add --no-cache memcached \
    && memcached -d -u memcache

COPY cmd.sh cmd.sh

RUN chmod a+x cmd.sh

CMD ["./cmd.sh"]
