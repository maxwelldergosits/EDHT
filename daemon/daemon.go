package main

import (
  "fmt"
  "flag"
  "os"
  "EDHT/utils"
  "EDHT/common/rpc_stubs"
  . "EDHT/common"
  "time"
  "github.com/mad293/mlog"
  "strings"
)
var (
  port string
  ip   string
  groupPort string
  groupAddress string
  dataDir string
  id uint64

  vall bool
  verboseLevels []string

  ml mlog.MLog

  data utils.StringStringMap
  preCommits utils.StringStringMap
)



func registerCLA(){

  flag.StringVar(&port, "port", "1456","Port to bind the server to")
  flag.StringVar(&ip, "address", "127.0.0.1","address to bind the server to")
  flag.StringVar(&dataDir, "data-dir", time.Now().String()+"/","data directory, please end with a /")

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

  preCommits = utils.NewStringStringMap(dataDir+"preCommits")
  data = utils.NewStringStringMap(dataDir+"data")
  ml = mlog.Create(verboseLevels,"",true,vall)

  ml.VPrintln("info","port:",port)
  ml.VPrintln("info","ip-address:",ip)


  rs := RemoteServer{
    Address:groupAddress,
    Port:groupPort,
    ID:0,
    Coordinator:true}

  rr,_,err := rpc_stubs.AttachToGroupRPC(false,ip,port,utils.GenMachineId(),rs)
  if (err!=nil) {
    ml.NPrintln("Couldn't join group, Exiting")
    os.Exit(1)
  }
  id =rr.ID
  startServer(ip,port)

}
