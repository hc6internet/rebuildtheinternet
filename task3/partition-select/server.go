package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
)

type msg struct {
	siteId      string
	partitionId string
}

var partitionTotal uint32 = 1

func main() {
	port := os.Getenv("PORT")
	pc := os.Getenv("PARTITION_COUNT")

	if sh, err := strconv.ParseUint(pc, 10, 32); err != nil {
		log.Printf("No partition count")
		return
	} else {
		partitionTotal = uint32(sh)
	}

	// zmq socket
	sock, err := zmq.NewSocket(zmq.REP)
	defer sock.Close()
	err = sock.Bind("tcp://*:" + port)
	if err != nil {
		log.Fatal(err)
		return
	}

	ch := make(chan msg)
	go processor(ch)

	for {
		handleMessage(sock, ch)
	}
}

func processor(ch chan msg) {
	h := fnv.New32a()

	for {
		m := <-ch

		h.Write([]byte(m.siteId))
		idx := h.Sum32() % partitionTotal

		var b strings.Builder
		fmt.Fprintf(&b, "%s%d", "db", idx)
		m.partitionId = b.String()

		ch <- m
	}
}

func handleMessage(sock *zmq.Socket, ch chan msg) {
	// read request
	bytes, err := sock.RecvBytes(0)
	if err != nil {
		log.Printf("Error recv: %s", err.Error())
		return
	}

	req := &PartitionRequest{}
	if err := proto.Unmarshal(bytes, req); err == nil {
		var m msg
		m.siteId = req.SiteId
		ch <- m
		m = <-ch

		// send response back
		resp := &PartitionResponse{PartitionId: m.partitionId, Seq: req.Seq, Status: true}
		out, err := proto.Marshal(resp)
		if err != nil {
			log.Printf("Error marshal response: ", err.Error())
			return
		}
		ret, err := sock.SendBytes(out, 0)
		if err != nil {
			log.Printf("Error writing response: ", err.Error())
			return
		}
		log.Printf("sent %d bytes", ret)
	}
}
