package main

import (
	"fmt"

	"net"
)

type Subscriber struct {
	Address net.Addr

	Connection net.Conn

	do chan Task
}

type Task struct {
	name string
}

var SubscriberCount int

var Subscribers []Subscriber

var CurrentSubscriber int

var taskChannel chan Task

func (sb Subscriber) awaitTask() {

	select {

	case t := <-sb.do:

		fmt.Println(t.name, "assigned")

	}

}

func serverListen(listener net.Listener) {

	for {

		conn, _ := listener.Accept()

		SubscriberCount++

		subscriber := Subscriber{Address: conn.RemoteAddr(), Connection: conn}

		subscriber.do = make(chan Task)

		subscriber.awaitTask()

		_ = append(Subscribers, subscriber)

	}

}

func doTask() {

	for {

		select {

		case task := <-taskChannel:

			fmt.Println(task.name, "invoked")

			Subscribers[CurrentSubscriber].do <- task

			if (CurrentSubscriber + 1) > SubscriberCount {

				CurrentSubscriber = 0

			} else {

				CurrentSubscriber++

			}

		}

	}

}

func main() {

	destinationStatus := make(chan int)

	SubscriberCount = 0

	CurrentSubscriber = 0

	taskChannel = make(chan Task)

	listener, err := net.Listen("tcp", ":9000")

	if err != nil {

		fmt.Println("Could not start server!", err)

	}

	go serverListen(listener)

	go doTask()

	<-destinationStatus

}
