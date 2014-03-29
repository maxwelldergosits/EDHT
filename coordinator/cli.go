package main

import (
  "flag"
  "os"
  "fmt"
)





func registerCLA(){

  // local options
  flag.StringVar(&port, "port", "1456","Port to bind the server to")
  flag.StringVar(&ip, "address", "127.0.0.1","address to bind the server to")
  flag.BoolVar(&verbose, "verbose", false, "verbose output")

  // group connection options
  flag.BoolVar(&groupconnect, "connect-to-group", false, "connect to an existing group of coordinators")
  flag.StringVar(&groupAddress, "group-address", "", "Address of any node in a group to connect to")
  flag.StringVar(&groupPort, "group-port", "", "Port of that the node in the group is on")

  //group configuration options
  flag.IntVar(&nshards, "shards",1,"Number of \"shards\" of data")
  flag.IntVar(&failures,"failures",0,"Number of failures tolerated")

  // local file options
  flag.StringVar(&logDir, "log-dir","","Directory output for log files (default is the current directory) directory must exist")
  flag.StringVar(&dataDir, "data-dir","","Directory output for data files (default is the current directory) directory must exist")
  flag.BoolVar(&disableLog, "disable-log",true,"Disable log file output")


  flag.Parse()

  if groupconnect && (groupAddress == "" || groupPort == "") {

    fmt.Println("")
    fmt.Println("If connecting to a group you must specify a group port and group address")
    fmt.Println("")
    fmt.Println("Usage:")
    fmt.Println("")
    flag.PrintDefaults()

    os.Exit(1)
  }

}







