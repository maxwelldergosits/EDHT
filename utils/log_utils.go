package utils


import (
  "log"
  "time"
  "io"
  "os"
  "os/user"
  "strings"
  "path/filepath"
)

func GenLogger(verbose bool,prefix string) (printer func(a ...interface{}), verbosePrinter func(a ...interface{})) {


  dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
  if err != nil {
    log.Fatal(err)
  }
  if (prefix == "") {
    prefix = dir
  }
  usr, _ := user.Current()
  homeDir := usr.HomeDir
  if prefix[:1] == "~" {
    prefix = strings.Replace(prefix, "~", homeDir, 1)
  }

  var logFile *os.File

  var time = time.Now()

  const RFC3339 = "2006-01-02T15:04:05Z07:00"
  var logFileName = prefix+"/"+time.Format(RFC3339)+".log"
  log.Println("log file:",logFileName)

  logFile,_ = os.Create(logFileName)

  normalWriter := io.MultiWriter(logFile,os.Stdout)
  var verboseWriter io.Writer
  if (verbose) {
    verboseWriter = io.MultiWriter(logFile,os.Stdout)
  } else {
    verboseWriter = logFile
  }

  ver := func(a ... interface{}) {
    log.SetOutput(verboseWriter)
    log.Println(a)
  }

  norm := func(a ...interface{}) {
    log.SetOutput(normalWriter)
    log.Println(a)
  }

  return norm, ver

}
