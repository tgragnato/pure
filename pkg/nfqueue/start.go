package nfqueue

func Start(queueNum int, windowSizeMin uint, windowSizeMax uint) {
	for i := 0; i < queueNum; i++ {
		go worker(
			uint16(i),
			windowSizeMin,
			windowSizeMax,
			i >= queueNum/2,
		)
	}
}
