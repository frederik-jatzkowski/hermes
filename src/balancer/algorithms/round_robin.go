package algorithms

import (
	"sync"

	"fleo.software/infrastructure/hermes/server"
)

func NewRoundRobin(servers *[]server.Server) *RoundRobin {
	return &RoundRobin{
		servers: servers,
		length:  uint(len(*servers)),
		index:   0,
		mutex:   sync.Mutex{},
	}
}

type RoundRobin struct {
	servers *[]server.Server
	length  uint
	index   uint
	mutex   sync.Mutex
}

func (r *RoundRobin) Next() *server.Server {
	var s *server.Server // to hold the current server pointer
	var imod uint        // to hold the current server index
	r.mutex.Lock()       // enter critical section
	for i := 1; i <= len(*r.servers); i++ {
		// find first available server
		imod = (r.index + uint(i)) % r.length
		s = &(*r.servers)[imod]
		if s.CanHandle() {
			r.index = imod
			r.mutex.Unlock() // leave critical section
			return s
		}
	}
	r.mutex.Unlock() // leave critical section
	return nil
}
