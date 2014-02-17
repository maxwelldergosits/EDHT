package main

import (
  "log"
  "net/rpc"
  "net/http"
  "net"
  . "EDHT/common"
  "errors"
)


var (

  state State = State{0,map[int]NodeInfo{}}


)

type Member int

func (t *Member) Register(info * NodeInfo, res * int) error {

  if info.Address == "" {
    return errors.New("need address")
  }

  log.Println("new member registered, id:",state.NextID)
  log.Println("type:",info.MemberType)
  log.Println("ip:",info.Address)
  log.Println("port:",info.Port)

  info.Id = state.NextID
  info.Token = 0
  state.NextID+=1

  state.Nodes[info.Id] = *info
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
