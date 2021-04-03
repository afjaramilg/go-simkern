package main

import (
	"fmt"
	"pryct1/req"
    "pryct1/simk"
)


func reqTest() {
    r := req.Req{1, 2, 3, 4, 5}
	fmt.Printf("%s\n", r)

	buf := make([]byte, req.ReqBufSize)
    r2 := req.Req{}
    req.ReqDeserial(&r2, buf)
    fmt.Printf("%s\n", r)
}


func main() {
    simk.StartSimK()
    //tclient.StartTestClient()

}
