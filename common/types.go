package common



type NodeInfo struct {
  Port string
  Address string
  MemberType int
  Id int
  Token int
}

type State struct {
  NextID int
  Nodes map[int]NodeInfo
}


