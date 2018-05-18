package service

import (
	"net"
	"bufio"
)

type Point struct {

	X float64
	Y float64
	Ip int16

}

type Channel struct{

	connection net.Conn
	send chan Point

}


func NewChannel(conn net.Conn) *Channel {

	c := &Channel{
		connection:conn,
		send:make(chan Point,4),
	}

	return c
}

func (c *Channel) reader(){
	buf := bufio.NewReader(c.connection)
	for{
		pkg,_ := readPkg(buf)
		c.hanlde(pkg)
	}
}

func (c *Channel) writer(){
	buf := bufio.NewWriter(c.connection)
	for pkg := range c.send{
		_ := writePkg(buf)
		buf.Flush()
	}
}
