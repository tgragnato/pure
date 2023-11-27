package sni

import (
	"net"
)

func Listen(iface string) error {

	listener, err := net.Listen("tcp", iface+":443")
	if err != nil {
		return err
	}

	go func() {
		for {
			flow, err := listener.Accept()
			if err != nil {
				continue
			}
			go establishFlow(flow)
		}
	}()

	return nil
}
