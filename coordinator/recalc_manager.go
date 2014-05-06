package main
import (
    "time"
    "EDHT/utils"
    . "EDHT/common"
    )


func startRecalc(n int) {


  go loop(n)

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
keys64 := make([]uint64,len(keys))
          for i,v := range keys {keys64[i] = uint64(v)}

        if (utils.CV(keys64) < 1.0) {return}



        diffs,npts:= pts.CalculateDiffs(keys)
          ml.VPrintln("recalc","keys:",keys)
          ml.VPrintln("recalc", "ranges:",npts.Ranges())

          if !npts.Verify() {
            ml.VPrintln("recalc","new partition invalid")
          }
        ml.VPrintln("recalc","diffs:",diffs)

          rs := npts.Ranges()
          err = gc.UpdatePartitions(diffs,Ranges{rs}) //if it fails, thats okay

          if err != nil {
            ml.VPrintln("recalc", "Recalculation Error:",err.Error())
          } else {
            ml.VPrintln("recalc", "Successful Recalculation of keyspace")
          }

}

