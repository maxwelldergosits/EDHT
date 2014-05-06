package utils


func Mean(in []uint64) float64 {

  s := 0.0
  for _,v := range in {
    s += float64(v)
  }
  return (s/float64(len(in)))
}

func Dev(mean, v float64) float64{
  return ((mean - v)*(mean -v))
}

func Devs(keys []uint64) []float64 {
  mean := Mean(keys)
  out := make([]float64,len(keys))
  for i,iv := range keys {
    v := float64(iv)
    out[i] = Dev(mean,v)
  }
  return out
}

func StdDev(in []uint64) float64 {
  mean := Mean(in)
  s := 0.0
  for _,iv := range in {
    v := float64(iv)
    s += Dev(mean,v)
  }
  return s/float64(len(in))
}

func CV(in []uint64) float64 {

  mean := Mean(in)
  sd   := StdDev(in)

  return (sd/mean)

}
