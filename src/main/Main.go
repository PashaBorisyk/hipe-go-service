package main

import (
	"net/http"
	"golang.org/x/net/websocket"
)

func main()  {

	println("Point service starting...")

	http.Handle("/",websocket.Handler(WebSocketHandler))


}


func WebSocketHandler(ws *websocket.Conn){

}
