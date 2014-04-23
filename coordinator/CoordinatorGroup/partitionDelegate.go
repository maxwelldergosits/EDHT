package CoordinatorGroup
import (
  "EDHT/coordinator/CoordinatorGroup/group"
  "github.com/mad293/mlog"
  . "EDHT/common"
)
type PD struct {
  gms group.Group
  ml mlog.MLog
}

func (pd * PD) GetDaemon(id uint64) RemoteServer {
  return pd.gms.GetDaemon(id)
}

func (pd * PD) GetLocalID() uint64 {
  return pd.gms.GetID()
}

func (pd * PD) DeleteDaemon(id uint64) {
  pd.gms.Delete(pd.gms.GetDaemon(id))
}

func (pd * PD) Logger() *mlog.MLog {
  return &pd.ml
}
