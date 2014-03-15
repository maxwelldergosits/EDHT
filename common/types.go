package common

type RemoteServer struct{
  Address     string
  Port        string
  ID          int64
  Coordinator bool
}

type DaemonData struct{
  ServerDetails RemoteServer
  Hashtable     Hashtable
}

type RegisterReply struct {
  Coordinators     map[int64]RemoteServer
  Daemons          map[int64]RemoteServer
  Failures         int
  Shards           int
  ID               int64
}


type Tuple struct {
  Key   string
  Value string
}

type Hashtable struct {
	Size int
	Store map[string]string
}

