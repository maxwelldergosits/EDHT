package partition
import . "EDHT/common"

func init() {

  f = .05
  thres = .05
}

func pd(e,a uint) float32 {

  af := float32(a)
  ef := float32(e)

  if (af == ef) {
    return 0.0
  }
  return (ef - af)*2 / (ef + af)

}

var f float32
var thres float32

func (self *Shard) Copy() *Shard {
    r := new(Shard)
    *r = *self
    return r
}


func (t * PartitionSet) Copy() (PartitionSet) {

  var newShards []Shard = make([]Shard,len(t.Shards))
  for i := range t.Shards {
    oS := t.Shards[i]
    newMap := make(map[uint64]bool)
    for k,v := range oS.Daemons {
      newMap[k] =v
    }
    newShards[i] = Shard{oS.Start,oS.End,newMap,oS.Keys,oS.delegate}
  }
  return PartitionSet{newShards,t.d,false,Ranges{},0}
}

func (o* PartitionSet) Recalc(keys []uint) PartitionSet{
  t := o.Copy()

  for i:=1; i < len(t.Shards); i++ {

    shard_1 := t.Shards[i-1]
    shard_2 := t.Shards[i]

    s1 := shard_1.Start
    e1 := shard_1.End

    s2 := shard_2.Start
    e2 := shard_2.End

    k1 := keys[i-1]
    k2 := keys[i]

    pdiff := pd(k1,k2)
    o.d.Logger().VPrintln("recalc","pdiff=",pdiff)

    if pdiff < -thres {
      e_1n := uint64(-pdiff * f * float32(e2-s2) + float32(e1))
      s_2n := e_1n+1

      t.Shards[i].Start = s_2n
      t.Shards[i-1].End = e_1n

    }else if pdiff > thres {
      e_1n := uint64(-pdiff * f * float32(e1-s1) + float32(e1))
      s_2n := e_1n+1

      t.Shards[i].Start = s_2n
      t.Shards[i-1].End = e_1n
    }

  }
  return t

}

func (ps * PartitionSet) Verify() bool {
  for i:=1;i < len(ps.Shards); i++{
    if ps.Shards[i].Start != ps.Shards[i-1].End+1 {
      return false
    }
  }
  return true
}

func GenerateDiffs(oldPS, newPS PartitionSet) ([]Diff) {

  diffs := make([]Diff,0)

  for i := range oldPS.Shards {

    s := oldPS.Shards[i].Start
    e := oldPS.Shards[i].End
    sn := newPS.Shards[i].Start
    en := newPS.Shards[i].End
    if sn < s && i > 0 {
      //copy sn -> s from shard[i-1] to i
      diffs = append(diffs, Diff{i-1,i,sn,s})
    }
    if en > e && i < len(oldPS.Shards)-1 {
      //copy e -> en from shard[i+1]
      diffs = append(diffs, Diff{i+1,i,e,en})
    }
    if sn > s {
      // delete sn -> s after done copying
      diffs = append(diffs, Diff{-1,i,sn,s})
    }
    if en < e {
      // delete en -> e after done copying
      diffs = append(diffs, Diff{-1,i,en,e})
    }
  }
  return diffs

}

func (ps * PartitionSet) Ranges() []Range {
  ranges := make([]Range,len(ps.Shards))
  for i, v:= range ps.Shards {
    ranges[i] = Range{unconv(v.Start),unconv(v.End)}
  }
  return ranges
}
func (ps * PartitionSet) IntRanges() []uint64 {
  ranges := make([]uint64,2*len(ps.Shards))
  for i, v:= range ps.Shards {
    ranges[2*i] = v.Start
    ranges[2*i+1] = v.End
  }
  return ranges
}

func (ps * PartitionSet) CalculateDiffs(keys []uint) ([]Diff,*PartitionSet) {
  cp := ps.Recalc(keys)
  if(!cp.Verify()){panic("new partition is missing keys")}
  diffs := GenerateDiffs(*ps,cp)
  return diffs,&cp
}
