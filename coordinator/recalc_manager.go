package main
import (
  "time"
)


func startRecalc() {


  go loop(10)

}

func loop(n int) {

  for ;; {
    recalc()
    time.Sleep(time.Duration(n)* time.Second)
  }

}

func recalc() {
  ml.VPrintln("recalc", "Starting Recalculation of keyspace")
  pts := gc.GetPartitions()

  // for shard get number of keys being held
  keys,err := pts.GetNKeysForEachShard()

  if err != nil {
    ml.VPrintln("recalc", "Recalculation Error:",err.Error())
    return
  }

  diffs,newPTS:= pts.CalculateDiffs(keys)

  err = gc.UpdatePartitions(diffs,newPTS) //if it fails, thats okay
  if err != nil {
    ml.VPrintln("recalc", "Recalculation Error:",err.Error())
  } else {
    ml.VPrintln("recalc", "Successful Recalculation of keyspace")
  }

}

