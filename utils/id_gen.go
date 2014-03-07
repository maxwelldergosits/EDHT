package utils

<<<<<<< HEAD

import (
  "encoding/binary"
  "bytes"
  "net"
  "time"
)
func GenId(machineid int64, coor bool) int64{

  timestamp := time.Now().Unix()
  if coor {
    return timestamp << 22 | machineid << 1 | 1
  } else {
    return timestamp << 22 | machineid << 1 | 0
  }

}

// returns mac address of local computer
func GenMachineId() int64{

  interfaces,_ := net.Interfaces()
  buf := bytes.NewBuffer(interfaces[0].HardwareAddr) // b is []byte
  myfirstint, _ := binary.ReadVarint(buf)

  return myfirstint
}
