package types

type Node struct {
	ID                string            `json:"id"`
	ListenAddr        string            `json:"listenAddr"`
	AdvertiseAddr     string            `json:"advertiseAddr"`
	RaftListenAddr    string            `json:"raftListenAddr"`
	RaftAdvertiseAddr string            `json:"raftAdvertiseAddr"`
	Status            string            `json:"status"`
	Labels            map[string]string `json:"labels"`
	Role              NodeRole          `json:"role"`
	RaftID            uint64            `json:"raftid"`
}

func (n Node) IsManager() bool {
	return n.Role == RoleManager || n.Role == RoleMix
}

func (n Node) IsAgent() bool {
	return n.Role == RoleAgent || n.Role == RoleMix
}

type NodeRole string

const (
	RoleManager NodeRole = "manager"
	RoleAgent   NodeRole = "agent"
	RoleMix     NodeRole = "mixed"
)
