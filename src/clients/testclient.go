package clients

import (
	"bufio"
	"fmt"
	"os"
	"pryct1/req"
	"strings"
	//"strconv"

)



type testClient struct {
	clientState
}

func printProcList(plist []byte) {
	var cind uint32
	var ctype uint8
	fmt.Println("-----------------------")
	for i := 0; i <= len(plist)-5; i += 5 {
		req.DeserU32(&cind, plist[i:])
		ctype = uint8(plist[i+4])

		fmt.Printf("PROC: %d, %s\n", cind, req.CtypeMap[ctype])
	}
	fmt.Println("------------------------")
}

func (t *testClient) recvLoop() {
	for {
		read, err := t.serverConn.Read(t.reqbuf)
		if err != nil || read == 0 {
			fmt.Println("cant reach server, bye")
            os.Exit(0) //not great
			return
		}

		//fmt.Println("ive dun read somethin")

		var recvdReq req.Req
		req.ReqDeserial(&recvdReq, t.reqbuf)
		fmt.Printf("recieved %d %s\n", read, recvdReq)

		if recvdReq.Plsz > 0 {
			plbuf := make([]byte, recvdReq.Plsz)
			read, err = t.serverConn.Read(plbuf)
			if recvdReq.Rtype == req.PRLIST {
				printProcList(plbuf)
			} else {
				fmt.Printf("message: %s\n", string(plbuf))
			}
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
	fmt.Println("enter process number")
	fmt.Scanf("%d\n", &appInd)

	resp := req.Req{
		Rtype: req.PRCLOSE,
		Src:   t.clientID,
		Info:  appInd,
	}

	req.ReqSerial(t.reqbuf, &resp)
	t.serverConn.Write(t.reqbuf)
}

func (t *testClient) listProcs() {
	resp := req.Req{
		Rtype: req.PRLIST,
		Src:   t.clientID,
	}

	req.ReqSerial(t.reqbuf, &resp)
	t.serverConn.Write(t.reqbuf)
}

func (t *testClient) openFM() {
	resp := req.Req{
		Rtype: req.FMOPEN,
		Src:   t.clientID,
	}

	req.ReqSerial(t.reqbuf, &resp)
	t.serverConn.Write(t.reqbuf)
}

func (t *testClient) passMSG() {
	var dst uint32
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("enter destination")
	fmt.Scanf("%d\n", &dst)

	fmt.Println("string to send:")
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimSpace(msg)

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
			strfn{"2. list open processes", tc.listProcs},
			strfn{"3. open file manager", tc.openFM},
			strfn{"4. send message to proc", tc.passMSG},
		}

		for _, s := range options {
			fmt.Printf("%s\n", s.str)
		}

		fmt.Scanf("%d\n", &comm)
		options[comm].fn()
	}
}
