package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

var xchgName string = "cache-hit"

var hits uint = 0
var miss uint = 0

func main() {
	port := os.Getenv("MON_PORT")

	// rabbitmq
	rmqHost := os.Getenv("RMQ_HOST")

	// set up rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@" + rmqHost + ":5672/")
	if err != nil {
		log.Println("Failed to connect to rabbitmq: " + err.Error())
		return
	} else {
		log.Println("Connected to rabbitmq")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	defer ch.Close()

	// only checking the cache-hit queue
	q, err := ch.QueueDeclare(
		xchgName, // name
		false,    // durable
		false,    // delete when usused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		xchgName, // exchange
		false,
		nil,
	)

	// check and display cache statistics
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	go func() {
		for d := range msgs {
			processMsg(d.Body)
		}
	}()

	http.HandleFunc("/", cacheStatHandler)

	log.Printf("Listen on port " + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func processMsg(data []byte) {
	log.Printf("Received a message: %s", string(data))
	if string(data) == "h" {
		hits++
	} else if string(data) == "m" {
		miss++
	}
}

func cacheStatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hits = %d\n", hits)
	fmt.Fprintf(w, "miss = %d\n", miss)
}
