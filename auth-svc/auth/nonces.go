package auth

import (
	"sync"
	"time"
)

const nonceTimeValidation = time.Second * 3

type nonces struct {
	sync.Mutex
	list     map[uint64]time.Time
	validate time.Duration
}

func newNonces(length int) *nonces {
	return &nonces{
		list:     make(map[uint64]time.Time, length),
		validate: nonceTimeValidation,
	}
}

func (n *nonces) Add(nonce uint64) {
	n.Lock()
	defer n.Unlock()
	n.list[nonce] = time.Now()
}

func (n *nonces) Remove(nonce uint64) {
	n.Lock()
	defer n.Unlock()
	delete(n.list, nonce)
}

func (n *nonces) In(nonce uint64) bool {
	n.Lock()
	defer n.Unlock()
	_, ok := n.list[nonce]
	return ok
}
