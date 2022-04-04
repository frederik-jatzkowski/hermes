package logging

import (
	"log"
	"strings"
	"sync"
	"time"
)

func NewLogCounter(message string, interval time.Duration, useIds bool) *LogCounter {
	res := &LogCounter{
		useIds:  useIds,
		mutex:   sync.Mutex{},
		message: message,
		ids:     make([]string, 0),
		count:   0,
	}
	go res.run(interval)
	return res
}

type LogCounter struct {
	useIds  bool
	mutex   sync.Mutex
	message string
	ids     []string
	count   uint
}

func (lc *LogCounter) run(interval time.Duration) {
	for {
		time.Sleep(interval)
		lc.mutex.Lock()
		if lc.count > 0 {
			if !lc.useIds {
				log.Printf(lc.message, lc.count)
			} else {
				log.Printf(lc.message, lc.count, strings.Join(lc.ids, ", "))
			}

			lc.count = 0
			lc.ids = lc.ids[:0]
		}
		lc.mutex.Unlock()
	}
}

func (lc *LogCounter) Increment() {
	lc.mutex.Lock()
	lc.count++
	lc.mutex.Unlock()
}

func (lc *LogCounter) IdentifierIncrement(id string) {
	lc.mutex.Lock()
	lc.count++
	if lc.useIds {
		for i := range lc.ids {
			if lc.ids[i] == id {
				lc.mutex.Unlock()
				return
			}
		}
		lc.ids = append(lc.ids, id)
	}
	lc.mutex.Unlock()
}
