package main

import (
  "log"
  "container/list"
  . "EDHT/utils"
)

var (
  localAddress string
  localPort    string
  remoteServers list.List
)


func InitLocalState(alocalAddress string, alocalPort string) {

  // validate IP and port
  if !ValidateIP(alocalAddress){
    log.Panic("Error: invalid IP address: ", alocalAddress)
  }


  if !ValidatePort(alocalPort){
    log.Panic("Error: invalid port: ", alocalPort)
  }

  localAddress = alocalAddress
  localPort    = alocalPort

  remoteServers = *list.New()

}



type RemoteServer struct{
  address string
  port string
}

func addToServers(rs * RemoteServer) {
  remoteServers.PushFront(*rs)
}

// this is an RPC method so we can remotely tell a coordinator to add a new remoteserver
func RegisterCoordinator(ns * RemoteServer, res * int) error{

  addToServers(ns)

  *res = 1
  return nil

}


