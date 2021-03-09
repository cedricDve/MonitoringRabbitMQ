package main

import (
	//dependecies
	"fmt" // to log
	log "github.com/sirupsen/logrus" // to log
	"github.com/streadway/amqp" // rabbit mq sdk
	"os"
)
//environment variable for PORTS
var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main(){
	consume()
}

func consume(){

	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" +rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")
	
	//handle connection failure
	if err !=nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	
	//create a channel to talk to the queue using the connection
	ch, err := conn.Channel()
	//handle the error
		if err !=nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	//declare the queue
	q, err := ch.QueueDeclare(
		"publisher", //name
		true, //durable
		false, // delete when unused
		false, //exclusive
		false, //no-wait
		nil, //arguments
	)
	//handle error
	if err !=nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}
	fmt.Println("Channel and Queue established!")
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(	
		q.Name, //queue
		"", //consumer
		false, //auto-ack
		false, //exclusive
		false, //no-local
		false, //no-wait
		nil, //args
	)

	if err !=nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	//listen channel, forever
	forever := make(chan bool)

	//go func, grabbing messag out of the body and send ack back to the queu
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			d.Ack(false)
		}
	}()

	fmt.Println("Running..")

	<-forever

}