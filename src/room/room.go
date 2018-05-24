package room

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"log"
	"model"
	"net/http"
)

type RoomChannel chan Point

type Point struct {
	x float64
	y float64
}

type Room struct {
	thisAddr          string
	serverAddr        string
	clientBuffSize    int
	serverBuffSize    int
	roomChanel        RoomChannel
	serverConnection  *websocket.Conn
	clientConnections map[uint]*websocket.Conn
}

func NewRoom(thisAddr, serverAddr string, clientBufSize, serverBufSize int, maxClientSize uint) *Room {

	return &Room{
		thisAddr:          thisAddr,
		serverAddr:        serverAddr,
		clientBuffSize:    clientBufSize,
		serverBuffSize:    serverBufSize,
		roomChanel:        make(chan Point, maxClientSize),
		clientConnections: make(map[uint]*websocket.Conn),
		serverConnection:  nil,
	}

}

func (room *Room) InitClientConnections() {
	log.Println("Initing client connection pull...")
	http.Handle("/", websocket.Handler(room.serveClientConnection))
	err := http.ListenAndServe(room.thisAddr, websocket.Handler(room.serveClientConnection))
	if err != nil {
		log.Print("Error while listening to connections")
		log.Fatal(err)
	}
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
	go room.serveServerConnection()
	log.Println("Connection created successfully")
	return nil
}

func (room *Room) serveServerConnection() {
	log.Print("Connection created with addres " + room.serverConnection.Request().RemoteAddr)

	room.serverConnection.MaxPayloadBytes = room.serverBuffSize
	room.serverConnection.Write([]byte("Handshake"))
	defer room.InitServerConnection()
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

	for i := 0; i < 10; i++ {
		raw, err := json.Marshal(model.NewModel())
		if err == nil {
			ws.Write(raw)
		}
	}

	log.Println("Connection " + ws.Request().RemoteAddr + " closed")

}
