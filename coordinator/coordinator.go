package main

import (
  "log"
  "flag"
)


const(
  Daemon = iota
  Coordinator
)



var (
  port string
  ip   string
  groupPort string
  groupAddress string
  verbose bool
  groupconnect bool
)


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
  log.Println("port:",port)
  log.Println("ip-address:",ip)
  log.Println("verbose:",verbose)
  startServer(ip,port)

}


