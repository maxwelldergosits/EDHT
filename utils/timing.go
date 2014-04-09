package utils

import "time"

var times map[string]time.Time

func init() {
  times = make(map[string]time.Time)

}

func Trace(s string) {
    times[s] = time.Now()
}

func Un(s string) int64{
    startTime := times[s]
    endTime := time.Now()
    return endTime.Sub(startTime).Nanoseconds()
}
