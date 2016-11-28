package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		go startThread(&wg, i)
		wg.Add(1)
	}
	wg.Wait()
}

func startThread(wg *sync.WaitGroup, i int) {
	defer wg.Done()
	fmt.Println(i)
}
