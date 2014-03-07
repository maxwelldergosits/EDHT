package common

type RemoteServer struct{
  Address     string
  Port        string
  ID          int64
  Coordinator bool
}

type Daemon struct{
  ServerDetails RemoteServer
  hashtable     Hashtable
}

type RegisterReply struct {
  Coordinators     map[int64]RemoteServer
  Daemons          map[int64]RemoteServer
  ID               int64
}


type Tuple struct {
        Key, Value string
}

type Hashtable struct {
	//can add metdata? (Karl)
	Size int
	Store map[string]string
}

