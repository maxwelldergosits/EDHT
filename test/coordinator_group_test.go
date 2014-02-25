package main

import (
  "testing"
)
func TestInitLocalState(t * testing.T) {

  var testAddress = "127.0.0.1"
  var testPort    = "1234"

  InitLocalState(testAddress,testPort)


  if localAddress != testAddress {
    t.Error("test address incorrect")
  }
  if localPort != testPort {
    t.Error("test port incorrect")
  }

}

func TestRegisterCoordinator(t * testing.T) {

  var rs = RemoteServer{"127.0.0.1","1234"}
  var res int

  RegisterCoordinator(&rs,&res)

  if res != 1 {
    t.Error("res incorrect")
  }

  var foundRs bool = false

  for e := remoteServers.Front(); e != nil; e = e.Next() {
    if e.Value == rs { foundRs = true; break; }
  }

  if foundRs != true {
    t.Error("didn't add rs")
  }

}
