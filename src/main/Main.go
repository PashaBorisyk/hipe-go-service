package main

import (
	"log"
	"room"
	"sync"
)

var wg = new(sync.WaitGroup)

func main() {

	log.Print("Starting application...")

	wg.Add(1)
	go configureServerConnections()

	log.Println("Application started")
	wg.Wait()
	log.Print("Program finished")

}

func configureServerConnections() {

	defer wg.Done()
	connectionsRoom := room.NewRoomFromConfig(wg)
	err := connectionsRoom.InitServerConnection()
	if err != nil {
		log.Println("Application will run in test mode!")
	}
	connectionsRoom.InitClientConnections()

}
