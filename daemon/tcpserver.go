/**
* TCPArithServer
 */

package main

import (
	"log"
  "net/http"
	"net/rpc"
	"net"
	"os"
)

func startServer(ip string, port string) {


  reg := new(Daemon)
  rpc.Register(reg)
  rpc.HandleHTTP()
  l,e := net.Listen("tcp",ip+":"+port)
  if e != nil {
    log.Fatal("listen error:", e)
  }
  http.Serve(l, nil)

}

func checkError(err error) {
	if err != nil {
		log.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
