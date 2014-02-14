package main

import "log"
import "net/rpc"


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


func main() {
  client, err := rpc.DialHTTP("tcp", "127.0.0.1:1456")
  if err != nil {
    log.Fatal("dialing:", err)
  }
  // Synchronous call
  var args = Args{"2456","127.0.0.1","test"}
  var reply = Return{"",0,0}
  err = client.Call("Member.Register", &args, &reply)
  if err != nil {
    log.Fatal("connect error:", err)
  }
  log.Println("id:",reply.Id)
}
