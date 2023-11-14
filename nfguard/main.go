package main

import (
	"flag"
	"sync"
)

var (
	windowSizeMin uint
	windowSizeMax uint
)

func main() {

	var queueNum int

	flag.IntVar(&queueNum, "queueNum", 8, "The number of NFQUEUEs to attach to")
	flag.UintVar(&windowSizeMin, "windowSizeMin", 60, "Minimum TCP Window")
	flag.UintVar(&windowSizeMax, "windowSizeMax", 90, "Maximum TCP Window")
	flag.Parse()

	var wg sync.WaitGroup
	for i := 0; i < queueNum; i++ {
		wg.Add(1)
		go queueWorker(uint16(i), i >= queueNum/2, &wg)
	}
	wg.Wait()
}
