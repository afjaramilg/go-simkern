package clients

import (
	"fmt"
	"math/rand"
	"net"
	"pryct1/req"
	"time"
)

const (
	clientWaitChance = 0.2
	clientErrChance  = 0.1
)


func init() {
    rand.Seed(time.Now().UnixNano())
}



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

	respIden := req.Req{
		Rtype: req.IDEN,
		Info:  ctype,
	}

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

	cs.clientID = respIden.Info

	return true
}

func (fc *clientState) sendReply(orig req.Req, rtype uint16, msg []byte) {
	replyReq := req.Req{
		Rtype: rtype,
		Src:   fc.clientID,
		Info:  orig.Src,
		Plsz:  uint32(len(msg)),
	}

	req.ReqSerial(fc.reqbuf, &replyReq)
	fc.serverConn.Write(fc.reqbuf)
	fc.serverConn.Write(msg)
}


func randWait() bool {
	chance := rand.Float64()
	errLim := clientWaitChance + clientErrChance

	if chance < clientWaitChance {
		wait := rand.Intn(5)
        fmt.Printf("waitng %d seconds\n", wait)
		time.Sleep(time.Duration(wait) * time.Second)

	} else if chance < errLim {
		return false

	}

	return true
}
