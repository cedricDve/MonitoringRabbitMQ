package main

import (//dependecies	
	"fmt" // to log
	 "net"
	"time"
	"net/http" // run http func
	"github.com/julienschmidt/httprouter"//http router -> web server for our API
	log "github.com/sirupsen/logrus" // to log
	"github.com/streadway/amqp" // rabbit mq sdk
	"os"
)

//environment variable for PORTS
var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main() {

	conn, err := amqp.DialConfig("amqp://" + rabbit_user + ":" +rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/", amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, time.Second*2)
		},
		Heartbeat:  time.Second, 
	})
	if err != nil {
		log.Fatalf("%s: %s", "Connection Failed", err)
	}
	log.Printf("conn: %v, err: %v", conn, err)
		

	//define webserver
	router := httprouter.New()

	// define a route => publish/ here is the message
	router.POST("/publish/:message", func (w http.ResponseWriter, r *http.Request, p httprouter.Params){
		//submit func
		if (conn.IsClosed()){ fmt.Println("connection closed") }
		submit(w, r, p)
	})

	fmt.Println("Running..")
	//startup the app with htpp and server on port 80
	log.Fatal(http.ListenAndServe(":80", router))

}
// now a webserver is up and running in go !
//create and define the submit fun

func submit(writer http.ResponseWriter, request *http.Request, p httprouter.Params) {
	
	//grab the message (POST request)
	message := p.ByName("message")

	fmt.Println("Received message is: " + message)

	//establish connection with the queu with std func of the rabbitmq sdk
	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" +rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")

	//handle connection failure
	
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}


	defer conn.Close() //defer the conection and close it when the function is exits

	//create a channel to talk to the queue using the connection
	ch, err := conn.Channel()

	//handle the error
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	//defer the channel, close it
	defer ch.Close()

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
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	// publish a message on the queue
	err = ch.Publish(
		"", //exchange
		q.Name, //routing key
		false, //mandatory
		false, //imediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		})

	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish a message on the queue", err)
	}
	
	fmt.Println("Successfully published")

}