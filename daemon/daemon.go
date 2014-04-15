package main

import (
  . "EDHT/common"
  "fmt"
  "EDHT/common/group"
  "flag"
  "os"
  "github.com/mad293/mlog"
  "strings"
)
var (
  port string
  ip   string
  groupPort string
  groupAddress string
  id uint64

  vall bool
  verboseLevels []string

  ml mlog.MLog

  data map[string]string
  preCommits map[string]string
)



func registerCLA(){

  flag.StringVar(&port, "port", "1456","Port to bind the server to")
  flag.StringVar(&ip, "address", "127.0.0.1","address to bind the server to")

  flag.BoolVar(&vall, "vall", false, "print all verbosity levels to stdout")
  var vl string
  flag.StringVar(&vl, "verbose", "", "levels of verbosity, : separated")

  flag.StringVar(&groupAddress, "group-address", "", "Address of any node in a group to connect to")
  flag.StringVar(&groupPort, "group-port", "", "Port of that the node in the group is on")


  flag.Parse()

  verboseLevels = strings.Split(vl,":")

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

  preCommits = make(map[string]string)
  data = make(map[string]string)
  ml = mlog.Create(verboseLevels,"",true,vall)

  ml.VPrintln("info","port:",port)
  ml.VPrintln("info","ip-address:",ip)

  group.InitGroup(ml,nil)

  id := group.JoinGroupAsDaemon(groupAddress,groupPort,ip,port)
  if (id == 0) {
    ml.NPrintln("Couldn't join group, Exiting")
    os.Exit(1)
  }
  startServer(ip,port)

}
