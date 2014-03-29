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
  . "EDHT/utils"
  "EDHT/web_interface"
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
  verbose bool

  logDir string
  dataDir string

  verboseLog func(a ...interface{})
  normalLog func(a ...interface{})
)



func main() {

  registerCLA()
  normalLog,verboseLog = GenLogger(verbose,logDir,disableLog)

  normalLog("coordinator starting up")

  verboseLog("port:",port)
  verboseLog("ip-address:",ip)

  group.InitGroup(verboseLog,normalLog,newDaemon)

  if(groupconnect) {
    g := group.JoinGroupAsCoordinator(groupAddress,groupPort,ip,port)
    MakeKeySpace(int(g.Nshards))
  } else {

      normalLog("creating group")
      normalLog("waiting for",failures +1, "coordinators")
      normalLog("waiting for",nshards * (failures +1), "daemons")
      MakeKeySpace(nshards)

      group.CreateGroup(ip,port,uint(nshards),uint(failures))
  }

  go web_interface.StartUp(verboseLog,port+"8")

  group.CoordinatorStartServer(ip,port)

}


