package utils


import (
  "encoding/binary"
  "bytes"
  "net"
  "time"
)
func GenId(machineid uint64, coor bool) uint64{

  timestamp := uint64(time.Now().UnixNano())
  if coor {
    return timestamp << 22 | machineid << 1 | 1
  } else {
    return timestamp << 22 | machineid << 1 | 0
  }

}

// returns mac address of local computer
func GenMachineId() uint64{

  en0,_ := net.InterfaceByName("en0")
  buf := bytes.NewBuffer(en0.HardwareAddr) // b is []byte
  myfirstint, _ := binary.ReadVarint(buf)

  return uint64(myfirstint)
}

