package sntp

import (
	"fmt"
	"net"
)

func Listen() error {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 123})
	if err != nil {
		return err
	}

	go func() {
		for {
			request := make([]byte, 512)
			rlen, remote, err := listener.ReadFromUDP(request[0:])
			if err != nil {
				continue
			}
			if rlen > 0 && validFormat(request) {
				go func() {
					if _, err := listener.WriteTo(generate(request), remote); err != nil {
						fmt.Println(err.Error())
					}
				}()
			}
		}
	}()

	return nil
}
