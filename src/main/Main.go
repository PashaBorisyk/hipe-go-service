package main

import (
	"log"
	"room"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	log.Print("Starting application...")

	wg.Add(1)
	go configureServerConnections(&wg)

	log.Println("Application started")
	wg.Wait()
	log.Print("Program finished")

}

func configureServerConnections(group *sync.WaitGroup) {

	defer group.Done()
	connectionsRoom := room.NewRoomFromConfig(group)
	err := connectionsRoom.InitServerConnection()
	if err != nil {
		log.Println("Application will run in test mode!")
	}
	connectionsRoom.InitClientConnections()

}
