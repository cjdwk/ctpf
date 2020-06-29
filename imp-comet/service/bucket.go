package service

import (
	"sync"
	"sync/atomic"

	"github.com/oofpgDLD/ctpf/imp-comet/api"
	"github.com/oofpgDLD/ctpf/imp-comet/conf"
)

// Bucket is a channel holder.
type Bucket struct {
	c     *conf.Bucket
	cLock sync.RWMutex        // protect the channels for chs
	chs   map[string]*Channel // map sub key to a channel
	// Group
	groups       map[string]*Group // bucket Group channels
	routines    []chan *api.BroadcastReq
	routinesNum uint64

	ipCnts map[string]int32
}

// NewBucket new a bucket struct. store the key with im channel.
func NewBucket(c *conf.Bucket) (b *Bucket) {
	b = new(Bucket)
	b.chs = make(map[string]*Channel, c.Channel)
	b.ipCnts = make(map[string]int32)
	b.c = c
	b.groups = make(map[string]*Group, c.Group)
	b.routines = make([]chan *api.BroadcastReq, c.RoutineAmount)
	for i := uint64(0); i < c.RoutineAmount; i++ {
		c := make(chan *api.BroadcastReq, c.RoutineSize)
		b.routines[i] = c
		go b.Groupproc(c)
	}
	return
}

// ChannelCount channel count in the bucket
func (b *Bucket) ChannelCount() int {
	return len(b.chs)
}

// GroupCount Group count in the bucket
func (b *Bucket) GroupCount() int {
	return len(b.groups)
}

// GroupsCount get all Group id where online number > 0.
func (b *Bucket) GroupsCount() (res map[string]int32) {
	var (
		groupID string
		group   *Group
	)
	b.cLock.RLock()
	res = make(map[string]int32)
	for groupID, group = range b.groups {
		if group.Online > 0 {
			res[groupID] = group.Online
		}
	}
	b.cLock.RUnlock()
	return
}

// Put put a channel according with sub key.
func (b *Bucket) Put(gid string, ch *Channel) (err error) {
	var (
		group *Group
		ok   bool
	)
	b.cLock.Lock()
	// close old channel
	if dch := b.chs[ch.Key]; dch != nil {
		dch.Close()
	}
	b.chs[ch.Key] = ch
	if gid != "" {
		if group, ok = b.groups[gid]; !ok {
			group = NewGroup(gid)
			b.groups[gid] = group
		}
	}
	b.ipCnts[ch.IP]++
	b.cLock.Unlock()
	if group != nil {
		err = group.Put(ch)
	}
	return
}

// Del delete the channel by sub key.
func (b *Bucket) Del(dch *Channel) {
	var (
		ok   bool
		ch   *Channel
		group *Group
	)
	b.cLock.Lock()
	if ch, ok = b.chs[dch.Key]; ok {
		if ch == dch {
			delete(b.chs, ch.Key)
		}
		// ip counter
		if b.ipCnts[ch.IP] > 1 {
			b.ipCnts[ch.IP]--
		} else {
			delete(b.ipCnts, ch.IP)
		}
	}
	b.cLock.Unlock()
	if group != nil && group.Del(ch) {
		// if empty Group, must delete from bucket
		b.DelGroup(group)
	}
}

// Channel get a channel by sub key.
func (b *Bucket) Channel(key string) (ch *Channel) {
	b.cLock.RLock()
	ch = b.chs[key]
	b.cLock.RUnlock()
	return
}

// Group get a Group by Groupid.
func (b *Bucket) Group(rid string) (Group *Group) {
	b.cLock.RLock()
	Group = b.groups[rid]
	b.cLock.RUnlock()
	return
}

// DelGroup delete a Group by Groupid.
func (b *Bucket) DelGroup(group *Group) {
	b.cLock.Lock()
	delete(b.groups, group.ID)
	b.cLock.Unlock()
	group.Close()
}

// BroadcastGroup broadcast a message to specified Group
func (b *Bucket) BroadcastGroup(arg *api.BroadcastReq) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.c.RoutineAmount
	b.routines[num] <- arg
}

// Groups get all Group id where online number > 0.
func (b *Bucket) Groups() (res map[string]struct{}) {
	var (
		groupID string
		group   *Group
	)
	res = make(map[string]struct{})
	b.cLock.RLock()
	for groupID, group = range b.groups {
		if group.Online > 0 {
			res[groupID] = struct{}{}
		}
	}
	b.cLock.RUnlock()
	return
}

// IPCount get ip count.
func (b *Bucket) IPCount() (res map[string]struct{}) {
	var (
		ip string
	)
	b.cLock.RLock()
	res = make(map[string]struct{}, len(b.ipCnts))
	for ip = range b.ipCnts {
		res[ip] = struct{}{}
	}
	b.cLock.RUnlock()
	return
}

// UpGroupsCount update all Group count
func (b *Bucket) UpGroupsCount(GroupCountMap map[string]int32) {
	var (
		groupID string
		group   *Group
	)
	b.cLock.RLock()
	for groupID, group = range b.groups {
		group.AllOnline = GroupCountMap[groupID]
	}
	b.cLock.RUnlock()
}

// Groupproc
func (b *Bucket) Groupproc(c chan *api.BroadcastReq) {
	for {
		arg := <-c
		if g := b.Group(arg.GroupId); g != nil {
			g.Push(arg.Proto)
		}
	}
}
