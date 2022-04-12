package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"sync"
	"testing"
)

func StressTest1000(t *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	m := sync.Mutex{}
	success := 0
	dialFail := 0
	writeFail := 0
	readFail := 0
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := tls.Dial("tcp", "hermes.fleo.software:443", nil)
			if err != nil {
				m.Lock()
				dialFail++
				m.Unlock()
				return
			}
			_, err = conn.Write([]byte("GET /"))
			if err != nil {
				m.Lock()
				writeFail++
				m.Unlock()
				return
			}
			res, err := io.ReadAll(conn)
			fmt.Println(len(res))
			fmt.Print(string(res))
			if err != nil {
				m.Lock()
				readFail++
				m.Unlock()
				return
			}
			m.Lock()
			success++
			m.Unlock()
			defer wg.Done()
		}()
	}
	fmt.Printf("successful: %v, dial failed: %v, write failed: %v, res failed: %v \n", success, dialFail, writeFail, readFail)
	wg.Wait()
}
