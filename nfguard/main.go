package main

import (
	"flag"
	"sync"
)

func main() {

	var (
		queueNum      int
		windowSizeMin uint
		windowSizeMax uint
	)

	flag.IntVar(&queueNum, "queueNum", 8, "The number of NFQUEUEs to attach to")
	flag.UintVar(&windowSizeMin, "windowSizeMin", 60, "Minimum TCP Window")
	flag.UintVar(&windowSizeMax, "windowSizeMax", 90, "Maximum TCP Window")
	flag.Parse()

	go analytics.Report()

	var wg sync.WaitGroup
	for i := 0; i < queueNum; i++ {
		wg.Add(1)
		go queueWorker(
			uint16(i),
			windowSizeMin,
			windowSizeMax,
			i >= queueNum/2,
			&wg,
		)
	}
	wg.Wait()
}
