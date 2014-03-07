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

  addToServers(&rs)

  var foundRs bool = false

  for _,v := range remoteServers {
    if v == rs {
      foundRs = true
    }
  }

  if foundRs != true {
    t.Error("didn't add rs")
  }

}


