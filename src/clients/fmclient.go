package clients

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"pryct1/req"
	"strconv"
	"strings"
)

const (
	prefixPath = "fakefs/"
	maxLogMsgs = 100
)

type fmClient struct {
	clientState
	logFile *os.File
	logPos  []int64
}

func formatPrlist(plbuf []byte) string {
	retstr := ""
	var cind uint32
	var ctype uint8

	for i := 0; i <= len(plbuf)-5; i += 5 {
		req.DeserU32(&cind, plbuf[i:])
		ctype = uint8(plbuf[i+4])
		retstr += fmt.Sprintf("%d %s, ", cind, req.CtypeMap[ctype])
	}

	return retstr

}

func (fc *fmClient) parseLog(logbuf []byte) {
	log.SetOutput(fc.logFile)

	var origReq req.Req
	req.ReqDeserial(&origReq, logbuf)

	payload := ""
	if origReq.Plsz > 0 {
		switch origReq.Rtype {
		case req.PRLIST:
			payload = formatPrlist(logbuf[req.ReqBufSize:])
		default:
			payload = string(logbuf[req.ReqBufSize:])
		}
	}

	logstrPos, _ := fc.logFile.Seek(0, 1) //write pos before writting
	fc.logPos = append(fc.logPos, logstrPos)

	payload = strings.TrimSpace(payload)
	logstr := fmt.Sprintf("REQ: %s, PAYLOAD: %s", origReq, payload)
	log.Printf("%s", logstr)

}

func (fc *fmClient) sendLog(recvdReq req.Req, msg string) {
	lastN, _ := strconv.Atoi(msg)
    fmt.Printf("last %d logs!", lastN)

	startPos := int64(0)
	if lastN < len(fc.logPos) {
		logInd := len(fc.logPos) - lastN
		startPos = fc.logPos[logInd]
	}


	fc.logFile.Seek(startPos, 0)
	reader := bufio.NewReader(fc.logFile)

	line, err := reader.ReadBytes('\n')
	for err == nil {
		fc.sendReply(recvdReq, req.LOGMSG, line)
		line, err = reader.ReadBytes('\n')
	}

	fc.logFile.Seek(0, 2)
}

func (fc *fmClient) runCMD(recvdReq req.Req, cmd *exec.Cmd) {
	err := cmd.Run()
	var resp []byte
	if err != nil {
		resp = []byte(fmt.Sprintf("FM ERROR: error running %s", cmd))
		fc.sendReply(recvdReq, req.ERR, []byte(resp))
	} else {
		resp = []byte(fmt.Sprintf("FM OK: ran command %s", cmd))
		fc.sendReply(recvdReq, req.OK, []byte(resp))
	}
}

func (fc *fmClient) parseMsg(recvdReq req.Req, plbuf []byte) {
	msg := string(plbuf)
	fmt.Printf("command: %s\n", msg)

	fmcmd := strings.ToUpper(msg[:2])

	switch fmcmd {
	case "CR":
		dirnm := prefixPath + strings.TrimSpace(msg[2:])
		cmd := exec.Command("mkdir", dirnm)
		fc.runCMD(recvdReq, cmd)

	case "RM":
		dirnm := prefixPath + strings.TrimSpace(msg[2:])
		cmd := exec.Command("rm", "-r", dirnm)
		fc.runCMD(recvdReq, cmd)

	case "LG":
		fc.sendLog(recvdReq, strings.TrimSpace(msg[2:]))

	default:
		fc.sendReply(recvdReq, req.ERR,
			[]byte("FM ERROR: cant parse command"))
	}

}

func StartFMClient() {
	lfile, _ := os.OpenFile(prefixPath+"logFile",
		os.O_RDWR|os.O_CREATE, 0755)

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

		wait := randWait()
		if !wait {
			errMsg := []byte("FM ERROR: error")
			fc.sendReply(recvdReq, req.ERR, errMsg)
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
