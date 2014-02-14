package main

import (
  "fmt"
  "flag"
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
  fmt.Println("port:",port)
  fmt.Println("ip-address:",ip)
  fmt.Println("verbose:",verbose)
  startServer(ip,port)

}


