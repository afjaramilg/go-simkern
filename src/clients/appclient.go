package clients

import (
	"fmt"
	"pryct1/req"
)

type appClient struct {
	clientState
}

func StartAppClient() {
	ac := appClient{clientState{
		srvrAddr: "localhost",
		srvrPort: "8000",
		reqbuf:   make([]byte, req.ReqBufSize),
	}}

	ok := ac.serverConnect(req.PROC)
	if !ok {
		return
	}

	fmt.Printf("connected w/ clientID %d\n", ac.clientID)
	var recvdReq req.Req
	for {
		read, err := ac.serverConn.Read(ac.reqbuf)
		if err != nil || read == 0 {
			fmt.Println("cant reach server, bye")
			return
		}

		//fmt.Println("ive dun read somethin")

		req.ReqDeserial(&recvdReq, ac.reqbuf)
		fmt.Printf("recieved %d %s\n", read, recvdReq)

		if recvdReq.Plsz > 0 {
			plbuf := make([]byte, recvdReq.Plsz)
			read, err = ac.serverConn.Read(plbuf)
			fmt.Printf("message: %s\n", string(plbuf))
		}

		if wait := randWait(); !wait {
			msg := fmt.Sprintf("APP %d ERROR", ac.clientID)
			ac.sendReply(recvdReq, req.ERR, []byte(msg))
			continue
		}

		msg := fmt.Sprintf("APP %d RECIEVED MSG", ac.clientID)
		ac.sendReply(recvdReq, req.OK, []byte(msg))

		if recvdReq.Rtype == req.PRCLOSE {
			fmt.Println("good nite")
			return
		}

	}
}
