/*

Usage:





*/
package main

import (
  "log"
  "flag"
  . "EDHT/utils"
  . "EDHT/common"
  "EDHT/web_interface"
)



var (
  port string
  ip   string
  groupPort string
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

func registerCLA(){

  flag.StringVar(&port, "port", "1456","Port to bind the server to")
  flag.StringVar(&ip, "address", "127.0.0.1","address to bind the server to")
  flag.BoolVar(&verbose, "verbose", false, "verbose output")
  flag.BoolVar(&groupconnect, "connect-to-group", false, "connect to an existing group of coordinators")
  flag.StringVar(&groupAddress, "group-address", "", "Address of any node in a group to connect to")
  flag.StringVar(&groupPort, "group-port", "", "Port of that the node in the group is on")
  flag.StringVar(&logDir, "log-dir","","Directory output for log files (default is the current directory) directory must exist")
  flag.StringVar(&dataDir, "data-dir","","Directory output for data files (default is the current directory) directory must exist")


  flag.Parse()
}

func main() {

  registerCLA()
  normalLog,verboseLog = GenLogger(verbose,logDir)

  normalLog("coordinator starting up")

  verboseLog("port:",port)
  verboseLog("ip-address:",ip)

  InitLocalState(ip,port,groupconnect)

  if(groupconnect) {

    AttachToGroup(groupAddress,groupPort)

  }
  go web_interface.StartUp()

  startServer(ip,port)

}


