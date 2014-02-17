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
  verbose bool
)


func registerCLA(){

  flag.StringVar(&port, "port", "1456","Port to bind the server to")
  flag.StringVar(&ip, "address", "127.0.0.1","address to bind the server to")
  flag.BoolVar(&verbose, "verbose", false, "verbose output")


  flag.Parse()

}

func main() {

  registerCLA()
  log.Println("port:",port)
  log.Println("ip-address:",ip)
  log.Println("verbose:",verbose)
  startServer(ip,port)

}


