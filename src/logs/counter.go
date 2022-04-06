package logs

import (
	"sync"
	"time"
)

var counters map[string]*LogCounter = make(map[string]*LogCounter)

func Counter(name string) *LogCounter {
	c, ok := counters[name]
	if !ok {
		c = &LogCounter{
			mutex: sync.Mutex{},
			name:  name,
			count: 0,
		}
		counters[name] = c
		go c.run()
	}
	return c
}

type LogCounter struct {
	mutex sync.Mutex
	name  string
	count int
}

func (c *LogCounter) run() {
	for {
		time.Sleep(time.Minute)
		c.mutex.Lock()
		if c.count > 0 {
			continuousLogger.Printf(c.name+": %v per minute", c.count)
			c.count = 0
		}
		c.mutex.Unlock()
	}
}

func (c *LogCounter) Increment() {
	c.mutex.Lock()
	c.count++
	c.mutex.Unlock()
}
