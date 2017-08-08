package discovery

type Server struct {
	Id      string
	Name    string
	Address string
	Port    int
	Extends map[string]interface{}
}

type ServerChange int8

const (
	SC_Online  ServerChange = iota + 1
	SC_Offline 
)

func (s ServerChange) String() string {
	switch s {
	case SC_Offline:
		return "Offline"
	case SC_Online:
		return "Online"
	}
	return ""
}

type ServersChanged func(server *Server, change ServerChange)

type Discoverer interface {
	GetServers() (servers map[string]*Server, err error)
	// 通知服务改变(下线,上线). 在第一次watch会一次性收到当前所有在线的服务上线通知
	WatchServer(ServersChanged)
	// 刷新存活时间
	UpdateServerTTL(leaseId string) (err error)
	RegisterService(server *Server) (leaseId string, err error)
	UnRegisterService(serverId string) (err error)
}
