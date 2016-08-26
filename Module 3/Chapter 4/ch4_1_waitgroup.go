package main

import(
	"time"
	"fmt"
	"sync"
)



var currentTime time.Time
var rwLock sync.RWMutex
func updateTime() {
	rwLock.RLock()
	currentTime = time.Now();
	time.Sleep(5 * time.Second)
	rwLock.RUnlock()
}

func main() {

	var wg sync.WaitGroup

	currentTime = time.Now();
	timer := time.NewTicker(2 * time.Second)
	writeTimer := time.NewTicker(10 * time.Second)
	endTimer := make(chan bool)

	wg.Add(1)
	go func() {

		for {
			select {
				case <- timer.C:
					fmt.Println(currentTime.String())
				case <- writeTimer.C:
					updateTime()
				case <- endTimer:
					timer.Stop()
					return
			}

		}

	}()

	wg.Wait()
	fmt.Println(currentTime.String())
}