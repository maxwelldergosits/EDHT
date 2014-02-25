package main

import (
  "log"
  "net/rpc"
  "net/http"
  "net"
)


type Member int

func startServer(ip string, port string) {

  reg := new(Member)
  rpc.Register(reg)
  rpc.HandleHTTP()
  l,e := net.Listen("tcp",ip+":"+port)
  if e != nil {
    log.Fatal("listen error:", e)
  }
  http.Serve(l, nil)

}
