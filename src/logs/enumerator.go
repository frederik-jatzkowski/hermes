package logs

import (
	"strings"
	"sync"
	"time"
)

var enums map[string]*LogEnumerator = make(map[string]*LogEnumerator)

func Enumerator(name string) *LogEnumerator {
	e, ok := enums[name]
	if !ok {
		e = &LogEnumerator{
			mutex:    sync.Mutex{},
			name:     name,
			messages: make([]string, 0),
		}
		enums[name] = e
		go e.run()
	}
	return e
}

type LogEnumerator struct {
	mutex    sync.Mutex
	name     string
	messages []string
}

func (e *LogEnumerator) run() {
	for {
		time.Sleep(time.Minute)
		e.mutex.Lock()
		if len(e.messages) > 0 {
			continuousLogger.Printf(e.name+" within the last minute: %v", strings.Join(e.messages, ", "))
			e.messages = e.messages[:0]
		}
		e.mutex.Unlock()
	}
}

func (e *LogEnumerator) Add(msg string) {
	e.mutex.Lock()
	for i := range e.messages {
		if msg == e.messages[i] {
			e.mutex.Unlock()
			return
		}
	}
	e.messages = append(e.messages, msg)
	e.mutex.Unlock()
}
