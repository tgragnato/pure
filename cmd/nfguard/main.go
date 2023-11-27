package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tgragnato/pure/pkg/nfqueue"
	"github.com/tgragnato/pure/pkg/shsh"
	"github.com/tgragnato/pure/pkg/sni"
	"github.com/tgragnato/pure/pkg/sntp"
)

func main() {

	var (
		iface4        string
		iface6        string
		queueNum      int
		windowSizeMin uint
		windowSizeMax uint
	)

	flag.StringVar(&iface4, "interfaceIPv4", "127.0.0.1", "Set here the IPv4 of the interface to bind to")
	flag.StringVar(&iface6, "interfaceIPv6", "[::1]", "Set here the IPv6 of the interface to bind to")
	flag.IntVar(&queueNum, "queueNum", 8, "The number of NFQUEUEs to attach to")
	flag.UintVar(&windowSizeMin, "windowSizeMin", 60, "Minimum TCP Window")
	flag.UintVar(&windowSizeMax, "windowSizeMax", 90, "Maximum TCP Window")
	flag.Parse()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	err := sntp.Listen()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}

	shsh.Listen(iface4, iface6)

	err = sni.Listen(iface4)
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}
	err = sni.Listen(iface6)
	if err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}

	nfqueue.Start(8, 60, 90)

	<-signalCh
}
