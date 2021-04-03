package clients

import (
	"fmt"
	"pryct1/req"
	//"strconv"
)

type testClient struct {
	clientState
}

func (t *testClient) recvLoop() {
	for {
		read, err := t.serverConn.Read(t.reqbuf)
		if err != nil || read == 0 {
			fmt.Println("cant reach server, bye")
			return
		}

		fmt.Println("ive dun read somethin")

		var recvdReq req.Req
		req.ReqDeserial(&recvdReq, t.reqbuf)
		fmt.Printf("recieved %d %s\n", read, recvdReq)

		if recvdReq.Plsz > 0 {
			plbuf := make([]byte, recvdReq.Plsz)
			read, err = t.serverConn.Read(plbuf)
			fmt.Printf("message: %s\n", string(plbuf))
		}
	}

}

func (t *testClient) openApp() {
	resp := req.Req{
		Rtype: req.PROPEN,
		Src:   t.clientID,
	}

	req.ReqSerial(t.reqbuf, &resp)
	t.serverConn.Write(t.reqbuf)
}

func (t *testClient) closeApp() {
	var appInd uint32
	fmt.Println("enter app number")
	fmt.Scanf("%d\n", &appInd)

	resp := req.Req{
		Rtype: req.PRCLOSE,
		Src:   t.clientID,
		Info:  appInd,
	}

	req.ReqSerial(t.reqbuf, &resp)
	t.serverConn.Write(t.reqbuf)
}

func (t *testClient) listApp() {
	resp := req.Req{
		Rtype: req.PRLIST,
		Src:   t.clientID,
	}

	req.ReqSerial(t.reqbuf, &resp)
	t.serverConn.Write(t.reqbuf)
}

func (t *testClient) openFM() {
	resp := req.Req{
		Rtype: req.PROPEN,
		Src:   t.clientID,
		Info:  1, // 1 is FM's clientID
	}

	req.ReqSerial(t.reqbuf, &resp)
	t.serverConn.Write(t.reqbuf)
}

func (t *testClient) closeFM() {
	resp := req.Req{
		Rtype: req.PRCLOSE,
		Src:   t.clientID,
		Info:  1, // 1 is FM's clientID
	}

	req.ReqSerial(t.reqbuf, &resp)
	t.serverConn.Write(t.reqbuf)
}

func (t *testClient) checkFM() {
	resp := req.Req{
		Rtype: req.FMCHECK,
		Src:   t.clientID,
	}

	req.ReqSerial(t.reqbuf, &resp)
	t.serverConn.Write(t.reqbuf)
}

func (t *testClient) passMSG() {
	var msg string
	var dst uint32

	fmt.Println("enter destination")
	fmt.Scanf("%d\n", &dst)

	fmt.Println("string 2 send:")
	fmt.Scanf("%s\n", &msg)

	resp := req.Req{
		Rtype: req.FWDMSG,
		Src:   t.clientID,
		Info:  dst,
		Plsz:  uint32(len(msg)),
	}

	req.ReqSerial(t.reqbuf, &resp)
	t.serverConn.Write(t.reqbuf)
	t.serverConn.Write([]byte(msg))
}

func StartTestClient() {
	tc := testClient{clientState{
		srvrAddr: "localhost",
		srvrPort: "8000",
		reqbuf:   make([]byte, req.ReqBufSize),
	}}

	ok := tc.serverConnect(req.USER)
	if !ok {
		return
	}

	fmt.Printf("connected w/ clientID  %d\n", tc.clientID)
	go tc.recvLoop()

	for {
		var comm uint8
		type strfn struct {
			str string
			fn  func()
		}

		options := [...]strfn{
			strfn{"0. open application", tc.openApp},
			strfn{"1. close application", tc.closeApp},
			strfn{"2. list open applications", tc.listApp},
			strfn{"3. open file manager", tc.openFM},
			strfn{"4. close file manager", tc.closeFM},
			strfn{"5. is file manager active?", tc.checkFM},
			strfn{"6. pass a message to a program", tc.passMSG},
		}

		for _, s := range options {
			fmt.Printf("%s\n", s.str)
		}

		fmt.Scanf("%d\n", &comm)
		options[comm].fn()
	}
}
