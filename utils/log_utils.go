package utils


import (
  "log"
)




func GenLogger(verbose bool, file string) (printer func(a ...interface{}), verbosePrinter func(a ...interface{})) {
  var ver func(a ...interface{})
  if verbose {
    ver = func(a ...interface{}) {
      log.Println(a)

    }
  } else {
    ver = func(a ...interface{}) { }
  }

  var norm = func(a ...interface{}) { log.Println(a)}
  return norm, ver

}
