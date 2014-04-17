package main
import (
  "time"
)


func startRecalc() {


  go loop()

}

func loop() {

  for ;; {
    recalc()
    time.Sleep(3 * time.Second)
  }

}

func recalc() {
  ml.VPrintln("recalc", "Starting Recalculation of keyspace")
  pts := gc.GetPartitions()

  pts.UpdateInfo()

  diffs := pts.CalculateDiffs()

  err := gc.UpdatePartitions(diffs) //if it fails, thats okay
  if err != nil {
    ml.VPrintln("recalc", "Recalculation Error:",err.Error())
  } else {
    ml.VPrintln("recalc", "Successful Recalculation of keyspace")
  }

}

