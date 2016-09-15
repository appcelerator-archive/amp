package membership

import (
	"errors"
	"sync"

	"google.golang.org/grpc"

	"github.com/coreos/etcd/raft/raftpb"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/state/watch"
	"github.com/gogo/protobuf/proto"
)

var (
	// ErrIDExists is thrown when a node wants to join the existing cluster but its ID already exists
	ErrIDExists = errors.New("membership: can't add node to cluster, node id is a duplicate")
	// ErrIDRemoved is thrown when a node tries to perform an operation on an existing cluster but was removed
	ErrIDRemoved = errors.New("membership: node was removed during cluster lifetime")
	// ErrIDNotFound is thrown when we try an operation on a member that does not exist in the cluster list
	ErrIDNotFound = errors.New("membership: member not found in cluster list")
	// ErrConfigChangeInvalid is thrown when a configuration change we received looks invalid in form
	ErrConfigChangeInvalid = errors.New("membership: ConfChange type should be either AddNode, RemoveNode or UpdateNode")
	// ErrCannotUnmarshalConfig is thrown when a node cannot unmarshal a configuration change
	ErrCannotUnmarshalConfig = errors.New("membership: cannot unmarshal configuration change")
)

// Cluster represents a set of active
// raft Members
type Cluster struct {
	mu      sync.RWMutex
	members map[uint64]*Member

	// removed contains the list of removed Members,
	// those ids cannot be reused
	removed        map[uint64]bool
	heartbeatTicks int

	PeersBroadcast *watch.Queue
}

// Member represents a raft Cluster Member
type Member struct {
	*api.RaftMember

	api.RaftClient
	Conn   *grpc.ClientConn
	tick   int
	active bool
}

// NewCluster creates a new Cluster neighbors list for a raft Member.
// Member marked as inactive if there was no call ReportActive for heartbeatInterval.
func NewCluster(heartbeatTicks int) *Cluster {
	// TODO(abronan): generate Cluster ID for federation

	return &Cluster{
		members:        make(map[uint64]*Member),
		removed:        make(map[uint64]bool),
		heartbeatTicks: heartbeatTicks,
		PeersBroadcast: watch.NewQueue(),
	}
}

// Tick increases ticks for all members. After heartbeatTicks node marked as
// inactive.
func (c *Cluster) Tick() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, m := range c.members {
		if !m.active {
			continue
		}
		m.tick++
		if m.tick > c.heartbeatTicks {
			m.active = false
		}
	}
}

// Members returns the list of raft Members in the Cluster.
func (c *Cluster) Members() map[uint64]*Member {
	members := make(map[uint64]*Member)
	c.mu.RLock()
	for k, v := range c.members {
		members[k] = v
	}
	c.mu.RUnlock()
	return members
}

// Removed returns the list of raft Members removed from the Cluster.
func (c *Cluster) Removed() []uint64 {
	c.mu.RLock()
	removed := make([]uint64, 0, len(c.removed))
	for k := range c.removed {
		removed = append(removed, k)
	}
	c.mu.RUnlock()
	return removed
}

// GetMember returns informations on a given Member.
func (c *Cluster) GetMember(id uint64) *Member {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.members[id]
}

func (c *Cluster) broadcastUpdate() {
	peers := make([]*api.Peer, 0, len(c.members))
	for _, m := range c.members {
		peers = append(peers, &api.Peer{
			NodeID: m.NodeID,
			Addr:   m.Addr,
		})
	}
	c.PeersBroadcast.Publish(peers)
}

// AddMember adds a node to the Cluster Memberlist.
func (c *Cluster) AddMember(member *Member) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.removed[member.RaftID] {
		return ErrIDRemoved
	}
	member.active = true
	member.tick = 0

	c.members[member.RaftID] = member

	c.broadcastUpdate()
	return nil
}

// RemoveMember removes a node from the Cluster Memberlist, and adds it to
// the removed list.
func (c *Cluster) RemoveMember(id uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.removed[id] = true

	return c.clearMember(id)
}

// ClearMember removes a node from the Cluster Memberlist, but does NOT add it
// to the removed list.
func (c *Cluster) ClearMember(id uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.clearMember(id)
}

func (c *Cluster) clearMember(id uint64) error {
	m, ok := c.members[id]
	if ok {
		if m.Conn != nil {
			m.Conn.Close()
		}
		delete(c.members, id)
	}
	c.broadcastUpdate()
	return nil
}

// ReplaceMemberConnection replaces the member's GRPC connection and GRPC
// client.
func (c *Cluster) ReplaceMemberConnection(id uint64, oldConn *Member, newConn *Member) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	oldMember, ok := c.members[id]
	if !ok {
		return ErrIDNotFound
	}

	if oldConn.Conn != oldMember.Conn {
		// The connection was already replaced. Don't do it again.
		newConn.Conn.Close()
		return nil
	}

	oldMember.Conn.Close()

	newMember := *oldMember
	newMember.Conn = newConn.Conn
	newMember.RaftClient = newConn.RaftClient
	c.members[id] = &newMember

	return nil
}

// IsIDRemoved checks if a Member is in the remove set.
func (c *Cluster) IsIDRemoved(id uint64) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.removed[id]
}

// Clear resets the list of active Members and removed Members.
func (c *Cluster) Clear() {
	c.mu.Lock()
	for _, member := range c.members {
		if member.Conn != nil {
			member.Conn.Close()
		}
	}

	c.members = make(map[uint64]*Member)
	c.removed = make(map[uint64]bool)
	c.mu.Unlock()
}

// ReportActive reports that member is acive (called ProcessRaftMessage),
func (c *Cluster) ReportActive(id uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	m, ok := c.members[id]
	if !ok {
		return
	}
	m.tick = 0
	m.active = true
}

// Active returns true if node is active.
func (c *Cluster) Active(id uint64) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m, ok := c.members[id]
	if !ok {
		return false
	}
	return m.active
}

// ValidateConfigurationChange takes a proposed ConfChange and
// ensures that it is valid.
func (c *Cluster) ValidateConfigurationChange(cc raftpb.ConfChange) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.removed[cc.NodeID] {
		return ErrIDRemoved
	}
	switch cc.Type {
	case raftpb.ConfChangeAddNode:
		if c.members[cc.NodeID] != nil {
			return ErrIDExists
		}
	case raftpb.ConfChangeRemoveNode:
		if c.members[cc.NodeID] == nil {
			return ErrIDNotFound
		}
	case raftpb.ConfChangeUpdateNode:
		if c.members[cc.NodeID] == nil {
			return ErrIDNotFound
		}
	default:
		return ErrConfigChangeInvalid
	}
	m := &api.RaftMember{}
	if err := proto.Unmarshal(cc.Context, m); err != nil {
		return ErrCannotUnmarshalConfig
	}
	return nil
}

// CanRemoveMember checks if removing a Member would not result in a loss
// of quorum, this check is needed before submitting a configuration change
// that might block or harm the Cluster on Member recovery
func (c *Cluster) CanRemoveMember(from uint64, id uint64) bool {
	members := c.Members()
	nreachable := 0 // reachable managers after removal

	for _, m := range members {
		if m.RaftID == id {
			continue
		}

		// Local node from where the remove is issued
		if m.RaftID == from {
			nreachable++
			continue
		}

		if c.Active(m.RaftID) {
			nreachable++
		}
	}

	nquorum := (len(members)-1)/2 + 1
	if nreachable < nquorum {
		return false
	}

	return true
}
