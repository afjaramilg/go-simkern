package clients

import (
	"fmt"
	"net"
	"pryct1/req"
)

type clientState struct {
	srvrAddr, srvrPort string
	serverConn         net.Conn
	reqbuf             []byte
	clientID           uint32
}

func (cs *clientState) serverConnect(ctype uint32) bool {
	conn, err := net.Dial("tcp", cs.srvrAddr+":"+cs.srvrPort)
	if err != nil {
		fmt.Printf("error %s\n", err)
		return false
	}

    cs.serverConn = conn

	var respConn, respIden req.Req
	_, err = cs.serverConn.Read(cs.reqbuf)
	if err != nil {
		fmt.Println("error talking to server respConn")
		return false
	}

	req.ReqDeserial(&respConn, cs.reqbuf)
	fmt.Printf("CONN RESP: %s\n", respConn)

	if respConn.Rtype != req.OK {
		fmt.Println("error obtaining ID")
		return false
	}

	cs.clientID = respConn.Info

	respIden.Rtype = req.IDEN
	respIden.Src = cs.clientID
	respIden.Info = ctype

	req.ReqSerial(cs.reqbuf, &respIden)
	cs.serverConn.Write(cs.reqbuf)
	_, err = cs.serverConn.Read(cs.reqbuf)
    if err != nil {
		fmt.Println("error talking to server respIden")
		return false
	}

	req.ReqDeserial(&respIden, cs.reqbuf)
	fmt.Printf("IDEN RESP: %s\n", respIden)

	if respIden.Rtype != req.OK {
        fmt.Println("error identifying self")
        return false
	}

    return true
}
