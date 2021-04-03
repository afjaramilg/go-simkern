package simk

import (
	//"fmt"
	"net"
	//"sync"
)

type clientConn struct {
	Conn     net.Conn
	ClientID uint32
	Ctype    uint8
	//mut      *sync.Mutex
}

/*
func (c *clientConn) SendBuf(buf []byte) bool {
	fmt.Printf("hello there %d!\n", c.ClientID)
    c.mut.Lock()
    start := 0
	for start < len(buf) {
		written, err := c.Conn.Write(buf[start:])
		start += written

		if err != nil || written <= 0 {
			fmt.Println("there was an error sending buf")
			c.mut.Unlock()
			return false
		}
	}
	c.mut.Unlock()
	return true
}

func (c *clientConn) RecvBuf(buf []byte) bool {
	c.mut.Lock()
	start := 0
	for start < len(buf) {
	    read, err := c.Conn.Read(buf[start:])
		start += read
		if err != nil || read <= 0 {
			fmt.Println("there was an error receving buf")
			c.mut.Unlock()
			return false
		}
	}
   	c.mut.Unlock()
	return true
}
*/
