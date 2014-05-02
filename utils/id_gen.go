package utils


import (
  "encoding/binary"
  "bytes"
  "net"
  "time"
)

func GetTimeNano() uint64 {
  return uint64(time.Now().UnixNano())
}

func GenId(machineid uint64, coor bool) uint64{
  var ts uint64
  timestamp := uint64(time.Now().UnixNano())
  if coor {
    ts = timestamp << 22 | machineid << 1 | 1
  } else {
    ts = timestamp << 22 | machineid << 1 | 0
  }
  return ts
}

// returns mac address of local computer
func GenMachineId() uint64{

  en0, err := net.InterfaceByName("en0")
  if (err != nil) {
    en0, err = net.InterfaceByName("eth0")
    if (err != nil) {
      panic("couldn't find a network id")
    }
  }
  buf := bytes.NewBuffer(en0.HardwareAddr) // b is []byte
  myfirstint, _ := binary.ReadVarint(buf)

  return uint64(myfirstint)
}


