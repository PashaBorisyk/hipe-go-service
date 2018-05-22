package main

import (
	"log"
	"room"
)

const (
	ThisAddr      = ":8081"
	PayloadLength = 1024
	ServerAddr    = "ws://localhost:9090/ws/"
)

func main() {

	log.Print("Point service starting...")
	configureServerConnections()
	log.Print("Program finished")

}

func configureServerConnections() {

	connectionsRoom := room.NewRoom(ThisAddr, ServerAddr, PayloadLength, PayloadLength, 10)
	err := connectionsRoom.InitServerConnection()
	if err != nil {
		log.Println("Application will run in test mode!")
	}
	connectionsRoom.InitClientConnections()

}
