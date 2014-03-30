package utils


import (
  "net"
  . "EDHT/common"
  "strconv"
  "net/rpc"
)

func ValidateIP(ip string) bool {

   if net.ParseIP(ip) == nil {
    return false
  }
  return true

}

func ValidatePort(port string) bool {

  portInt,err := strconv.Atoi(port)
  if err != nil || portInt > 65535 || portInt < 1 {
    return false
  }
  return true

}

// right now doesn't do much right now will cache connections when I get around to it (Maxwell)
func MakeConnection(rs RemoteServer) (*rpc.Client, error) {

  client,err := rpc.DialHTTP("tcp",rs.Address+":"+rs.Port)

  return client,err


}
