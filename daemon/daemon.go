package main

import (
  . "EDHT/common"
  "fmt"
  "EDHT/common/group"
  "flag"
  "os"
  "mlog"
)
var (
  port string
  ip   string
  groupPort string
  id uint64
  groupAddress string
  verbose bool
  ml mlog.MLog

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
  ml = mlog.Create([]string{},"",true,verbose)

    ml.VPrintln("debug","port:",port)
    ml.VPrintln("debug","ip-address:",ip)

  InitTPC()
  group.InitGroup(ml,nil)

  id := group.JoinGroupAsDaemon(groupAddress,groupPort,ip,port)
  if (id == 0) {
    ml.NPrintln("Couldn't join group, Exiting")
    os.Exit(1)
  }
  startServer(ip,port)

}
