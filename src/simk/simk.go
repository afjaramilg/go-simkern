package simk

import (
	"fmt"
	"log"
	"net"
	"pryct1/req"
	"sync"
)

type clientReq struct {
	r  req.Req
	cl clientConn
}

var conns sync.Map //clientID, clientConn

var (
	enter_ch = make(chan net.Conn)
	state_ch = make(chan clientReq)
)

const maxOpenProcs = 10

func remClient(cind uint32) {
	c_raw, found := conns.LoadAndDelete(cind)
	if found {
		c := c_raw.(clientConn)
		c.Conn.Close()
	}
}

func sendMsg(c clientConn, rtype uint16, payload []byte) bool {
	reqbuf := make([]byte, req.ReqBufSize)
	resp := req.Req{
		Id:    0,
		Rtype: rtype,
		Src:   0,
		Info:  c.ClientID,
		Plsz:  uint32(len(payload)),
	}

	req.ReqSerial(reqbuf, &resp)
	_, err1 := c.Conn.Write(reqbuf)
	_, err2 := c.Conn.Write(payload)

	if err1 != nil || err2 != nil {
		fmt.Println("error sending message")
		return false
	}

	return true
}

func fwdMsg(c clientConn, r req.Req) {
	plbuf := make([]byte, r.Plsz)
	read, err := c.Conn.Read(plbuf)

	if err != nil || read == 0 {
		remClient(c.ClientID)
		return
	}

	dst_raw, dstinmap := conns.Load(r.Info)
	if dst_raw == nil || !dstinmap {
		msg := fmt.Sprintf("dst %d doesnt exist", r.Info)
		sendMsg(c, req.ERR, []byte(msg))
	} else {
		dst := dst_raw.(clientConn)
		sent := sendMsg(dst, req.FWDMSG, plbuf)
		if !sent {
			msg := fmt.Sprintf("cant send msg to %d", r.Info)
			sendMsg(c, req.ERR, []byte(msg))
		} else {
			msg := "great success"
			sendMsg(c, req.OK, []byte(msg))
		}
	}
}

func clientHandler(cind uint32) {
	fmt.Printf("started clientHandler %d\n", cind)
	reqbuf := make([]byte, req.ReqBufSize)

	for {
		c_raw, srcinmap := conns.Load(cind)
		if c_raw == nil || !srcinmap {
			remClient(cind)
			break
		}

		c := c_raw.(clientConn)
		read, err := c.Conn.Read(reqbuf)
		if err != nil || read == 0 {
			remClient(cind)
			break
		}

		var creq req.Req
		req.ReqDeserial(&creq, reqbuf)
		fmt.Printf("req recieved %s\n", creq)

		if creq.Rtype == req.IDEN {
			c.Ctype = uint8(creq.Info)
			conns.Store(cind, c)
			sendMsg(c, req.OK, []byte{})
			fmt.Printf("client %d said he's a %d\n",
				c.ClientID, creq.Info)

		} else {
			switch {
			case c.Ctype == req.UNKN:
				msg := "identify client type first"
				sendMsg(c, req.ERR, []byte(msg))

			case (creq.Rtype >= req.PROPEN &&
				creq.Rtype <= req.FMCHECK):
				state_ch <- clientReq{creq, c}

			case creq.Rtype == req.FWDMSG:
				fwdMsg(c, creq)
			}
		}
	}

	fmt.Printf("closing clientHandler %d\n", cind)
}

func openProc(openProcs *uint32) (uint16, []byte) {
	if *(openProcs) >= maxOpenProcs {
		msg := fmt.Sprintf("cant open anymore processes")
		return req.ERR, []byte(msg)
	}

	*openProcs++
	msg := fmt.Sprintf("process number %d started", *openProcs)
	return req.OK, []byte(msg)
}

func closeProc(openProcs *uint32, procID uint32) (uint16, []byte) {
	if procID < 2 || procID > 31 {
		msg := fmt.Sprintf("invalid process ID %d", procID)
		return req.ERR, []byte(msg)
	}

	*openProcs--
	fmt.Println("[call to remove client here]")

	msg := fmt.Sprintf("process %d closed", procID)
	return req.OK, []byte(msg)
}

func stateHandler() {
	var openProcs uint32 = 0

	for clreq := range state_ch {
		var succ uint16
		var msg []byte

		switch clreq.r.Rtype {
		case req.PROPEN:
			fmt.Println("opening process")
			succ, msg = openProc(&openProcs)
		case req.PRCLOSE:
			fmt.Println("closing process")
			succ, msg = closeProc(&openProcs, clreq.r.Info)
		}
		sendMsg(clreq.cl, succ, msg)
	}

	fmt.Println("closin stateHandler")
}

func handleEnter() {
	var ind uint32 = 2
	for nconn := range enter_ch {
		newcc := clientConn{Conn: nconn, ClientID: ind}
		conns.Store(ind, newcc)

		sendMsg(newcc, req.OK, []byte{})
		go clientHandler(ind)
		ind++
	}

	fmt.Println("closin handleEnter")
}

func StartSimK() {
	listenAddr := "localhost"
	listenPort := "8000"
	listener, err := net.Listen("tcp", listenAddr+":"+listenPort)
	if err != nil {
		log.Fatal(err)
	}

	go handleEnter()
	go stateHandler()

	fmt.Println("starting simk bb")
	for {
		client, err := listener.Accept()

		if err != nil {
			log.Print(err)
			continue
		}

		enter_ch <- client
	}

}
