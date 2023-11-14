package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/florianl/go-nfqueue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func queueWorker(queueNum uint16, ipv6 bool, wg *sync.WaitGroup) {
	defer wg.Done()

	nf, err := nfqueue.Open(&nfqueue.Config{
		NfQueue:      queueNum,
		MaxPacketLen: 0xFFFF,
		MaxQueueLen:  0xFF,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: time.Second,
	})
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer nf.Close()

	fn := func(a nfqueue.Attribute) int {
		id := *a.PacketID
		var packet gopacket.Packet
		if ipv6 {
			packet = gopacket.NewPacket(*a.Payload, layers.LayerTypeIPv6, gopacket.Default)
		} else {
			packet = gopacket.NewPacket(*a.Payload, layers.LayerTypeIPv4, gopacket.Default)
		}

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			log.Println("Non TCP packet" + packet.String())
			err := nf.SetVerdict(id, nfqueue.NfAccept)
			if err != nil {
				log.Println(err.Error())
			}
			return 0
		}

		tcp, _ := tcpLayer.(*layers.TCP)
		if !tcp.SYN && !tcp.ACK {
			log.Println("Non ( SYN | ACK ) packet" + packet.String())
			err := nf.SetVerdict(id, nfqueue.NfAccept)
			if err != nil {
				log.Println(err.Error())
			}
			return 0
		}

		windowSize := rand.Float64()*float64(windowSizeMax-windowSizeMin) + float64(windowSizeMin)
		packet.TransportLayer().(*layers.TCP).Window = uint16(windowSize)
		err := packet.TransportLayer().(*layers.TCP).SetNetworkLayerForChecksum(packet.NetworkLayer())
		if err != nil {
			log.Println(err.Error())
		}
		buffer := gopacket.NewSerializeBuffer()
		options := gopacket.SerializeOptions{
			ComputeChecksums: true,
			FixLengths:       true,
		}
		if err := gopacket.SerializePacket(buffer, options, packet); err != nil {
			log.Println(err.Error())
		}
		packetBytes := buffer.Bytes()
		err = nf.SetVerdictModPacket(id, nfqueue.NfAccept, packetBytes)
		if err != nil {
			log.Println(err.Error())
		}
		return 0
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = nf.RegisterWithErrorFunc(ctx, fn, func(e error) int {
		if e != nil {
			log.Fatalln(e.Error())
		}
		return 0
	})
	if err != nil {
		log.Println(err.Error())
	}
	<-ctx.Done()
}
