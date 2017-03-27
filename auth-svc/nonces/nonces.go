package nonces

import (
	"sync"
	"time"
)

const windowTimeSeconds = 600 //

// Nonces is a safe thread list of nonces, designed to use as a black list.
type Nonces struct {
	sync.Mutex
	list   map[uint64]int64 //map[nonce]expiration
	window int64
}

// New creates a new list of nonces.
func New(length int) *Nonces {
	return &Nonces{
		list:   make(map[uint64]int64, length),
		window: windowTimeSeconds,
	}
}

// Clean removes expired nonces from list.
func (n *Nonces) Clean(nonce uint64) {
	n.Lock()
	defer n.Unlock()

	expired := make([]uint64, 0)
	for nonce, expiration := range n.list {
		if time.Now().Unix() > expiration {
			expired = append(expired, nonce)
		}
	}

	for _, nonce := range expired {
		delete(n.list, nonce)
	}
}

// CheckAndAdd checks if nonce if valid and adds it.
func (n *Nonces) CheckAndAdd(nonce uint64) bool {
	n.Lock()
	defer n.Unlock()

	_, ok := n.list[nonce]
	if !ok {
		n.list[nonce] = time.Now().Unix() + n.window
	}

	valid := !ok
	return valid
}
