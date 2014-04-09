/*

Usage of coordinator:
  -address="127.0.0.1": address to bind the server to
  -connect-to-group=false: connect to an existing group of coordinators
  -data-dir="": Directory output for data files (default is the current directory) directory must exist
  -disable-log=true: Disable log file output
  -failures=0: Number of failures tolerated
  -group-address="": Address of any node in a group to connect to
  -group-port="": Port of that the node in the group is on
  -log-dir="": Directory output for log files (default is the current directory) directory must exist
  -port="1456": Port to bind the server to
  -shards=1: Number of "shards" of data
  -verbose=false: verbose output



*/
package main

import (
  "EDHT/common/group"
  "EDHT/web_interface"
  "EDHT/utils"
  "os"
  "mlog"
)



var (
  port string
  ip   string

  nshards int
  failures int

  disableLog bool
  groupPort string

  groupAddress string
  groupconnect bool
  verboseLevels []string
  all bool

  logDir string
  dataDir string

  ml mlog.MLog
)



func main() {

  registerCLA()
  ml = mlog.Create(verboseLevels,"",true,all)

  ml.NPrintln("coordinator starting up")

  ml.VPrintln("debug","port:",port)
  ml.VPrintln("debug","ip-address:",ip)

  group.InitGroup(ml,NewDaemon)

  if(groupconnect) {
    utils.Trace("connect")
    g := group.JoinGroupAsCoordinator(groupAddress,groupPort,ip,port)
    ml.VPrintln("time","connected in ",utils.Un("connect")/1e6,"milliseconds")
    if group.GetLocalID() == 0 {
      ml.NPrintln("Couldn't join group shutting down")
      os.Exit(1)
    }
    MakeKeySpace(int(g.Nshards))
  } else {

      ml.NPrintln("creating group")
      ml.NPrintln("waiting for",failures +1, "coordinators")
      ml.NPrintln("waiting for",nshards * (failures +1), "daemons")
      MakeKeySpace(nshards)

      group.CreateGroup(ip,port,uint(nshards),uint(failures))
  }

  go web_interface.StartUp(ml,port+"8",GetK,PutKV)

  group.CoordinatorStartServer(ip,port)

}


