package clients

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"pryct1/req"
	"strings"
)

const (
	prefixPath = "fakefs/"
	maxLogMsgs = 20
)

type fmClient struct {
	clientState
	logFile *os.File
}

func (fc *fmClient) sendReply(orig req.Req, rtype uint16, msg []byte) {
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

func (fc *fmClient) parseLog(logbuf []byte) {
	log.SetOutput(fc.logFile)

	var origReq req.Req
	req.ReqDeserial(&origReq, logbuf)

	payload := ""
	if origReq.Plsz > 0 {
		payload = string(logbuf[req.ReqBufSize:])
	}

	log.Printf("REQ: %s, PAYLOAD: %s", origReq, payload)
}

func (fc *fmClient) sendLog(recvdReq req.Req) {
	reader := bufio.NewReader(fc.logFile)

	line, err := reader.ReadBytes('\n')
	for l := 0; l < maxLogMsgs && err == nil; l++ {
		fc.sendReply(recvdReq, req.FWDMSG, line)
		line, err = reader.ReadBytes('\n')
	}
}

func (fc *fmClient) runCMD(recvdReq req.Req, cmd *exec.Cmd) {
	err := cmd.Run()
	var resp []byte
	if err != nil {
		resp = []byte(fmt.Sprintf("FM ERROR: error running %s", cmd))
		fc.sendReply(recvdReq, req.FWDMSG, []byte(resp))
	} else {
		resp = []byte(fmt.Sprintf("FM OK: ran command %s", cmd))
		fc.sendReply(recvdReq, req.FWDMSG, []byte(resp))
	}
}

func (fc *fmClient) parseMsg(recvdReq req.Req, plbuf []byte) {
	msg := string(plbuf)
	fmt.Printf("command: %s\n", msg)

	fmcmd := strings.ToUpper(msg[:2])
	dirnm := prefixPath + strings.TrimSpace(msg[2:])

	switch fmcmd {
	case "CR":
		cmd := exec.Command("mkdir", dirnm)
		fc.runCMD(recvdReq, cmd)

	case "RM":
		cmd := exec.Command("rm", "-r", dirnm)
		fc.runCMD(recvdReq, cmd)

	case "LG":
		fc.sendLog(recvdReq)

	default:
		fc.sendReply(recvdReq, req.FWDMSG,
			[]byte("FM ERROR: cant parse command"))
	}

}

func StartFMClient() {
	lfile, _ := os.Create(prefixPath + "logFile")
	fc := fmClient{
		clientState: clientState{
			srvrAddr: "localhost",
			srvrPort: "8000",
			reqbuf:   make([]byte, req.ReqBufSize),
		},
		logFile: lfile,
	}

	ok := fc.serverConnect(req.FM)
	if !ok {
		return
	}

	fc.clientID = 1

	fmt.Printf("connected w/ clientID %d\n", fc.clientID)
	var recvdReq req.Req
	for {
		read, err := fc.serverConn.Read(fc.reqbuf)
		if err != nil || read == 0 {
			fmt.Println("cant reach server, bye")
			return
		}

		//fmt.Println("ive dun read somethin")

		req.ReqDeserial(&recvdReq, fc.reqbuf)
		fmt.Printf("recieved %d %s\n", read, recvdReq)

		plbuf := make([]byte, recvdReq.Plsz)

		if recvdReq.Plsz > 0 {
			fc.serverConn.Read(plbuf)
			fmt.Printf("%s\n", string(plbuf))
		}

		switch recvdReq.Rtype {
		case req.PRCLOSE:
			fmt.Println("good nite")
			return

		case req.FWDMSG:
			fc.parseMsg(recvdReq, plbuf)

		case req.LOGMSG:
			fc.parseLog(plbuf)

		}

	}
}
