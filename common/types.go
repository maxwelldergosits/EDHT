package common

type NodeInfo struct {
  Port string
  Address string
  MemberType int
  Id int
  Token int
}

type DaemonInfo struct {
  KeySpaceStart int
  KeySpaceEnd int
  Address string
  Port string
  Id int
  Token int
}

type CoordinatorInfo struct {
  Address string
  Port string
  Id int
  Token int
}

type State struct {
  NextID int
  Token int
  Daemons map[int]DaemonInfo
  Coordinators map[int]CoordinatorInfo
}

type DaemonState struct {
  Address string
  Port string
  Id int
  Token int
}
