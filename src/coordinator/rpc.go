package main

import (
  "log"
  "net/rpc"
  "net/http"
  "net"
  "errors"
)



var (
  nextId int = 0
)


type Args struct {
  Port string
  Address string
  MemType string
}

type Return struct {
  Success string
  Id int
  Token int
}

type Member int

func (t *Member) Register(a * Args, r * Return) error {
  if a.Address == "" {
    return errors.New("need address")
  }
  log.Println("new member registered, id:",nextId)
  log.Println("type:",a.MemType)
  log.Println("ip:",a.Address)
  log.Println("port:",a.Port)

  r.Id = nextId
  r.Success = "Success"
  r.Token = 0
  nextId+=1

  return nil
}


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
