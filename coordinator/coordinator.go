/*

Usage:





*/
package main

import (
  "log"
  . "EDHT/utils"
  . "EDHT/common"
  "EDHT/web_interface"
)



var (
  port string
  ip   string
  shards int
  failures int
  groupPort string
  disableLog bool
  groupAddress string
  verbose bool
  logDir string
  dataDir string
  groupconnect bool
  verboseLog func(a ...interface{})
  normalLog func(a ...interface{})
)

func InitLocalState(alocalAddress string, alocalPort string,connect bool) {

  // validate IP and port
  if !ValidateIP(alocalAddress){
    log.Panic("Error: invalid IP address: ", alocalAddress)
  }


  if !ValidatePort(alocalPort){
    log.Panic("Error: invalid port: ", alocalPort)
  }

  localAddress = alocalAddress
  localPort    = alocalPort

  remoteCoordinators = map[int64]RemoteServer{}
  remoteDaemons = map[int64]RemoteServer{}
  pendingCommits = map[int64]RemoteServer{}


  if !connect {

    machine_id := GenMachineId()

    id = GenId(machine_id,true)
    remoteCoordinators[id] = RemoteServer{localAddress,localPort,id,true}

  }
}


func main() {

  registerCLA()
  normalLog,verboseLog = GenLogger(verbose,logDir,disableLog)

  normalLog("coordinator starting up")

  verboseLog("port:",port)
  verboseLog("ip-address:",ip)

  InitLocalState(ip,port,groupconnect)

  if(groupconnect) {

    AttachToGroup(groupAddress,groupPort)

  } else {

      normalLog("creating group")
      normalLog("waiting for",failures +1, "coordinators")
      normalLog("waiting for",shards * (failures +1), "daemons")

  }

  go web_interface.StartUp(verboseLog)

  startServer(ip,port)

}


