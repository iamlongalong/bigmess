package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	rmap "github.com/iamlongalong/bigmess/pkg/responsivemap"
)

func init() {

}

var (
	EventStateChange        = "changeState"
	EventClusterStateChange = "changeClusterState"
)

func NewCluster() *Cluster {
	h := &Handlers{}
	c := &Cluster{
		peers:    map[string]*Peer{},
		handlers: h,
	}

	state := rmap.NewResponsiveMap()
	// 监听map变更

	c.cinfo = state

	h.Register(EventStateChange, c.HandleStateChanges)
	h.Register(EventClusterStateChange, c.HandleClusterStateChanges)

	// create a key listener
	baseKeyListener := rmap.NewMapListener("", func(me rmap.MapEvent) {
		c.mu.Lock()
		defer c.mu.Unlock()

		for _, cli := range c.peers {
			cli.ChangeState(context.Background(), &ChangeStateMessage{
				Option: string(me.Option),
				Key:    me.Key,
				Val:    me.NewVal,
			})
		}
	})

	// watch key
	c.cinfo.Watch(baseKeyListener)

	return c
}

type Cluster struct {
	mu sync.Mutex

	cinfo *rmap.ResponsiveMap

	peers map[string]*Peer

	handlers *Handlers

	// 自己的订阅状态
	substates *rmap.ResponsiveMap
}

func (clu *Cluster) SubRoom(roomid string) {
	// clu.substates.
}

type ClusterConfig struct {
	Addrs []string
}

// 连上后操作
func (clu *Cluster) NewClusterPeer(c *Peer) {
	clu.mu.Lock()
	defer clu.mu.Unlock()

	if oc := clu.peers[c.ipaddr]; oc != nil {
		log.Printf("clinent %s has already in cluster", c.ipaddr)
		c.cancel()
		return
	}

	clu.peers[c.ipaddr] = c

}

func (clu *Cluster) SyncStates(c *Peer) {

}

func (clu *Cluster) FindClusters(cfg ClusterConfig) {

	for _, addr := range cfg.Addrs {
		go RetryTimes(func() error {
			if !clu.hasClusterNode(addr) { // 不存在 node，建立连接
				conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
				if err != nil {
					log.Printf("new connect fail : %s", err)
					return err
				}

				c := NewPeer(UUIDV4(), addr, conn, clu.handlers)
				c.Start()
				clu.NewClusterPeer(c)
			}

			return nil
		}, 5, time.Second*2)

	}
}

func (clu *Cluster) hasClusterNode(addr string) bool {
	for _, c := range clu.peers {
		if c.ipaddr == addr {
			return true
		}
	}

	return false
}
