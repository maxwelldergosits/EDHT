package main

import (
  . "EDHT/common"
  "fmt"
  "EDHT/common/group"
  "flag"
  "os"
  "EDHT/utils"
)
var (
  port string
  ip   string
  groupPort string
  id uint64
  groupAddress string
  verbose bool
  verboseLog func(a...interface{})
  normalLog func(a...interface{})

  local_state DaemonData
)


func registerCLA(){

  flag.StringVar(&port, "port", "1456","Port to bind the server to")
  flag.StringVar(&ip, "address", "127.0.0.1","address to bind the server to")

  flag.BoolVar(&verbose, "verbose", false, "verbose output")
  flag.StringVar(&groupAddress, "group-address", "", "Address of any node in a group to connect to")
  flag.StringVar(&groupPort, "group-port", "", "Port of that the node in the group is on")


  flag.Parse()

  if groupAddress == "" || groupPort == "" {
    fmt.Println("Usage:")
    fmt.Println("")
    fmt.Println("   Group Address and Group Port required")
    fmt.Println("")
    flag.PrintDefaults()
    os.Exit(1)
  }
}



func main() {

 registerCLA()

  local_state = SpawnDaemon(ip,port)
  normalLog,verboseLog = utils.GenLogger(verbose,"",true)

    verboseLog("port:",port)
    verboseLog("ip-address:",ip)

  InitTPC()
  group.InitGroup(verboseLog,normalLog,nil)

  id := group.JoinGroupAsDaemon(groupAddress,groupPort,ip,port)
  if (id == 0) {
    normalLog("Couldn't join group, Exiting")
    os.Exit(1)
  }
  startServer(ip,port)

}
