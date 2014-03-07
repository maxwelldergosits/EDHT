package main

import (
  "log"
  "flag"
  . "EDHT/utils"
  . "EDHT/common"
)



var (
  port string
  ip   string
  groupPort string
  groupAddress string
  verbose bool
  groupconnect bool
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

  remoteServers = map[string]RemoteServer{}

}

func registerCLA(){

  flag.StringVar(&port, "port", "1456","Port to bind the server to")
  flag.StringVar(&ip, "address", "127.0.0.1","address to bind the server to")
  flag.BoolVar(&verbose, "verbose", false, "verbose output")
  flag.BoolVar(&groupconnect, "connect-to-group", false, "connect to an existing group of coordinators")
  flag.StringVar(&groupAddress, "group-address", "", "Address of any node in a group to connect to")
  flag.StringVar(&groupPort, "group-port", "", "Port of that the node in the group is on")


  flag.Parse()
}

func main() {

  registerCLA()

  if(verbose) {
    log.Println("port:",port)
    log.Println("ip-address:",ip)
  }
  InitLocalState(ip,port)

  if(groupconnect) {

    AttachToGroup(groupAddress,groupPort)

  }

  startServer(ip,port)

}


