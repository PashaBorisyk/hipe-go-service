package room

import (
	"config"
	"encoding/json"
	"golang.org/x/net/websocket"
	"log"
	"model"
	"net/http"
	"sync"
)

type RoomChannel chan Point

type Point struct {
	x float64
	y float64
}

type Room struct {
	thisPort          string
	serverAddr        string
	clientBuffSize    int
	serverBuffSize    int
	roomChanel        RoomChannel
	serverConnection  *websocket.Conn
	clientConnections map[int]*websocket.Conn
	waitGroup         *sync.WaitGroup
}

func NewRoom(thisAddr, serverAddr string, clientBufSize, serverBufSize int, maxClientSize int, waitGroup *sync.WaitGroup) *Room {

	return &Room{
		thisPort:          thisAddr,
		serverAddr:        serverAddr,
		clientBuffSize:    clientBufSize,
		serverBuffSize:    serverBufSize,
		roomChanel:        make(chan Point, maxClientSize),
		clientConnections: make(map[int]*websocket.Conn),
		serverConnection:  nil,
		waitGroup:         waitGroup,
	}

}

func NewRoomFromConfig(waitGroup *sync.WaitGroup) *Room {

	var globalConfig = config.GetConfig()

	return &Room{
		thisPort:          globalConfig.ConnectionsConfig.Client.ListenPort,
		serverAddr:        globalConfig.ConnectionsConfig.Server.Url,
		clientBuffSize:    globalConfig.ConnectionsConfig.Client.MaxConnectionPoolSize,
		roomChanel:        make(chan Point, globalConfig.ConnectionsConfig.Client.MaxConnectionPoolSize),
		clientConnections: make(map[int]*websocket.Conn),
		serverConnection:  nil,
	}

}

func (room *Room) InitClientConnections() {
	log.Println("Initing client connection pull...")
	room.waitGroup.Add(1)
	http.Handle("/", websocket.Handler(room.serveClientConnection))
}

func (room *Room) InitServerConnection() error {
	log.Println("Creating server connection...")
	conn, err := websocket.Dial(room.serverAddr, "ws", room.serverAddr)
	if err != nil {
		log.Println("Unable to connect to server with addres " + room.serverAddr + " protocol : ws")
		log.Println(err)
		return err
	}
	room.serverConnection = conn
	room.waitGroup.Add(1)
	go room.serveServerConnection()
	log.Println("Connection created successfully")
	return nil
}

func (room *Room) serveServerConnection() {
	log.Print("Connection created with addres " + room.serverConnection.Request().RemoteAddr)

	defer room.waitGroup.Done()

	room.serverConnection.MaxPayloadBytes = room.serverBuffSize
	for {
		msg := make([]byte, room.clientBuffSize)
		_, err := room.serverConnection.Read(msg)
		if err != nil {
			log.Println("Closing server connection with addres " + room.serverConnection.Request().RemoteAddr)
			log.Println(err)
			room.serverConnection.Close()
			return
		}

	}

}

func (room *Room) serveClientConnection(ws *websocket.Conn) {

	log.Print("Connection " + ws.Request().RemoteAddr + " created")
	ws.MaxPayloadBytes = room.clientBuffSize

	defer room.waitGroup.Done()

	for {
		raw, err := json.Marshal(model.NewModel())
		if err == nil {
			ws.Write(raw)
		}
		read, err := ws.Read(raw)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(read)
	}

	log.Println("Connection " + ws.Request().RemoteAddr + " closed")

}
