package req

import (
	"fmt"
	"strings"
)

const ( //rtype enum
	OK      = iota //ok
	ERR            //err
	IDEN           //identify self as ctype
	PROPEN         //open program
	PRCLOSE        //close program
	PRLIST         //list open programs
	FMOPEN         // open file manager
	FWDMSG         //pass a message to a program
	LOGMSG         //log somethin
)

var RtypeMap = map[uint16]string{
	OK:      "OK",
	ERR:     "ERR",
	IDEN:    "IDEN",
	PROPEN:  "PROPEN",
	PRCLOSE: "PRCLOSE",
	PRLIST:  "PRLIST",
	FMOPEN:  "FMOPEN",
	FWDMSG:  "FWDMSG",
	LOGMSG:  "LOGMSG",
}

const ( //Ctype enum
	UNKN = iota
	USER
	PROC
	FM
)

var CtypeMap = map[uint8]string{
	UNKN: "UNKN",
	USER: "USER",
	PROC: "PROC",
	FM:   "FM",
}

// opts is unsued for now
type Req struct {
	Id, Rtype uint16 // 2 * 2 bytes
	Src, Info uint32 // 2 * 4 bytes
	Plsz      uint32 // 4 bytes
} // total info size = 4 + 8 + 4 = 16 bytes

const ReqBufSize = 16

func (r Req) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("id: %d, ", r.Id))
	sb.WriteString(fmt.Sprintf("type: %s, ", RtypeMap[r.Rtype]))
	sb.WriteString(fmt.Sprintf("src: %d, info: %d, ", r.Src, r.Info))
	sb.WriteString(fmt.Sprintf("payload size: %d", r.Plsz))

	return sb.String()
}
