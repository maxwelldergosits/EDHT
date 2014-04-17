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
  "EDHT/web_interface"
  "EDHT/coordinator/CoordinatorGroup"
  "os"
  "github.com/mad293/mlog"
)


var (

  // final variables
  port string
  ip   string
  groupAddress string
  groupconnect bool
  logDir string
  dataDir string
  disableLog bool
  groupPort string
  ml mlog.MLog


  gc CoordinatorGroup.CoordinatorGroup // coordinated state.
)


func main() {

  verboseLevels, vall, nshards, failures:= registerCLA()
  ml = mlog.Create(verboseLevels,"",true,vall)

  ml.NPrintln("coordinator starting up")
  ml.VPrintln("debug","port:",port)
  ml.VPrintln("debug","ip-address:",ip)

  if(groupconnect) {

    gc, err := CoordinatorGroup.ConnectToGroup(groupAddress,groupPort,ip,port)

    if err != nil {
      ml.NPrintf("Error: %s, Couldn't join group shutting down\n",err.Error())
      os.Exit(1)
    }

  } else {

      ml.NPrintln("creating coordinator group")
      ml.NPrintln("waiting for",failures +1, "coordinators")
      ml.NPrintln("waiting for",nshards * (failures +1), "daemons")

      gc = CoordinatorGroup.NewCoodinatorGroup(nshards,failures,ml)
  }
  startRecalc()
  web_interface.StartUp(ml,port+"8",GetKey,PutKey)

}




