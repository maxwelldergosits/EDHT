package partition

func init() {

  f = .2
}

func pd(e,a uint) float32 {

  af := float32(a)
  ef := float32(e)

  return (ef - af)*2 / (ef + af)

}

var f float32

func (self *Shard) Copy() *Shard {
    r := new(Shard)
    *r = *self
    return r
}


func (t * PartitionSet) Copy() (PartitionSet) {

  var newShards []*Shard = make([]*Shard,len(t.shards))
  for i := range t.shards {
    oS := t.shards[i]
    newMap := make(map[uint64]bool)
    for k,v := range oS.daemons {
      newMap[k] =v
    }
    newShards[i] = &Shard{oS.Start,oS.End,newMap,oS.Keys}
  }
  return PartitionSet{newShards,t.d}
}

func (o* PartitionSet) Recalc() PartitionSet{
  t := o.Copy()

  for i:=1; i < len(t.shards); i++ {

    shard_1 := t.shards[i-1].Copy()
    shard_2 := t.shards[i].Copy()

    s1 := shard_1.Start
    e1 := shard_1.End

    s2 := shard_2.Start
    e2 := shard_2.End

    k1 := shard_1.Keys
    k2 := shard_2.Keys

    pdiff := pd(k1,k2)

    if pdiff < -.05 {
      e_1n := uint64(-pdiff * f * float32(e2-s2) + float32(e1))
      s_2n := e_1n+1

      t.shards[i].Start = s_2n
      t.shards[i-1].End = e_1n

    }else if pdiff > .05 {
      e_1n := uint64(-pdiff * f * float32(e1-s1) + float32(e1))
      s_2n := e_1n+1

      t.shards[i].Start = s_2n
      t.shards[i-1].End = e_1n
    }

  }
  return t

}


func GenerateDiffs(oldPS, newPS PartitionSet) ([]Diff,[]Diff) {

  toBeCopied := make([]Diff,len(oldPS.shards))
  toBeDeleted := make([]Diff,len(oldPS.shards))

  for i := range oldPS.shards {

    s := oldPS.shards[i].Start
    e := oldPS.shards[i].End
    sn := newPS.shards[i].Start
    en := newPS.shards[i].End
    if sn < s && i > 0 {
      //copy sn -> s from shard[i+1]
      toBeCopied = append(toBeCopied, Diff{i-1,i,sn,s})
    }
    if en > e && i < len(oldPS.shards)-1 {
      //copy e -> en from shard[i+1]
      toBeCopied = append(toBeCopied, Diff{i+1,i,e,en})

    }
    if sn > s {
      // delete sn -> s after done copying
      toBeDeleted = append(toBeDeleted, Diff{-1,i,sn,s})

    }
    if en < e {
      // delete en -> e after done copying
      toBeCopied = append(toBeDeleted, Diff{-1,i,e,en})
    }
  }
  return toBeCopied,toBeDeleted

}


func (t * PartitionSet) ApplyDiffs(cds []Diff,dds []Diff) {
  done := make(chan bool)
  for i := range cds {

    var diff = cds[i]
    go func() {
      to := t.shards[diff.To]
      from := t.shards[diff.From]
      start,end := diff.Start,diff.End

      done <- t.d.CopyDiff(to,from,start,end)
    }()

  }
  for i:=0; i<len(dds); i++ {
    <- done
  }
  for i := range dds {
    var diff = cds[i]
    go func() {
      from := t.shards[diff.From]
      start,end := diff.Start,diff.End

      t.d.DeleteDiff(from,start,end)
    }()
  }
}
