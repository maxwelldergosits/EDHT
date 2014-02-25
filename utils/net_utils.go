package utils

import (
  "net"
  "strconv"
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
