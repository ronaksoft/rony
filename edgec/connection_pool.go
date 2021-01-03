package edgec

import (
	"sync"
)

/*
   Creation Time: 2021 - Jan - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type connPool struct {
	mtx       sync.RWMutex
	pool      map[uint64]map[string]*wsConn
	leaderIDs map[uint64]string
}

func newConnPool() *connPool {
	cp := &connPool{
		pool:      make(map[uint64]map[string]*wsConn, 16),
		leaderIDs: make(map[uint64]string, 16),
	}
	return cp
}

func (cp *connPool) addConn(serverID string, replicaSet uint64, leader bool, c *wsConn) {
	cp.mtx.Lock()
	defer cp.mtx.Unlock()

	m := cp.pool[replicaSet]
	if m == nil {
		m = make(map[string]*wsConn, 16)
	}
	m[serverID] = c
	if leader {
		cp.leaderIDs[replicaSet] = serverID
	}
}

func (cp *connPool) removeConn(serverID string, replicaSet uint64) {
	cp.mtx.Lock()
	defer cp.mtx.Unlock()
}

func (cp *connPool) getConn(replicaSet uint64, onlyLeader bool) *wsConn {
	cp.mtx.RLock()
	defer cp.mtx.RUnlock()

	if onlyLeader {
		leaderID := cp.leaderIDs[replicaSet]
		if leaderID == "" {
			return nil
		}
		m := cp.pool[replicaSet]
		if m != nil {
			c := m[leaderID]
			go c.connect()
			return c
		}
	} else {
		m := cp.pool[replicaSet]
		if m != nil {
			for _, c := range m {
				go c.connect()
				return c
			}
		}
	}
	return nil
}
