package utils


import (
  "encoding/binary"
  "bytes"
  "net"
)
func GenId(timestamp int64, machineid int64, sequence_number int64 ) int64{

  return timestamp << 22 | machineid << 12 | sequence_number

}

// returns mac address of local computer
func GenMachineId() int64{

  interfaces,_ := net.Interfaces()
  buf := bytes.NewBuffer(interfaces[0].HardwareAddr) // b is []byte
  myfirstint, _ := binary.ReadVarint(buf)

  return myfirstint

}
