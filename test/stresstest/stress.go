package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	StressTest(n)
	end := time.Now()
	fmt.Printf("duration for %v requests: %v ms\n", n, end.Sub(start).Milliseconds())
}

func StressTest(n int) {
	wg := sync.WaitGroup{}
	m := sync.Mutex{}
	success := 0
	requestFail := 0
	responseFail := 0
	for i := 0; i < n; i++ {
		time.Sleep(time.Microsecond * 5000)
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := http.Get("https://hermes.fleo.software")
			if err != nil {
				m.Lock()
				requestFail++
				m.Unlock()
				return
			}
			// fmt.Println(len(res))
			b, err := io.ReadAll(res.Body)
			if err != nil || !strings.HasPrefix(string(b), "Lorem") {
				m.Lock()
				responseFail++
				m.Unlock()
				return
			}
			m.Lock()
			success++
			m.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("successful: %v, requests failed: %v, responses failed: %v\n", success, requestFail, responseFail)
}
