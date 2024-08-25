package sni

import (
	"crypto/rand"
	"crypto/tls"
	"net"
	"testing"
)

func Test_handleClientHello_Random(t *testing.T) {
	t.Parallel()

	randomData := make([]byte, 4096)
	_, err := rand.Read(randomData)
	if err != nil {
		t.Fatal(err)
	}

	client, server := net.Pipe()
	errChan := make(chan error)

	go func() {
		_, _, err := handleClientHello(server)
		errChan <- err
	}()

	go func() {
		_, err = client.Write(randomData)
		if err != nil {
			t.Error(err)
		}
		err = client.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	if err = <-errChan; err == nil {
		t.Error("expected error for random data")
	}
}

func Test_handleClientHello_HTTP(t *testing.T) {
	t.Parallel()

	client, server := net.Pipe()
	errChan := make(chan error)

	go func() {
		_, _, err := handleClientHello(server)
		errChan <- err
	}()

	go func() {
		_, err := client.Write([]byte("GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"))
		if err != nil {
			t.Error(err)
		}
		err = client.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	if err := <-errChan; err == nil {
		t.Error("expected error for HTTP/1.1 request")
	}
}

func Test_handleClientHello_SNI(t *testing.T) {
	t.Parallel()

	client, server := net.Pipe()
	clientHelloChan := make(chan *tls.ClientHelloInfo, 1)

	go func() {
		clientHelloInfo, _, err := handleClientHello(server)
		if err != nil {
			clientHelloChan <- nil
		}
		clientHelloChan <- clientHelloInfo
	}()

	go func() {
		tlsConn := tls.Client(
			client,
			&tls.Config{
				ServerName: "example.com",
				MinVersion: tls.VersionTLS12,
			},
		)
		if tlsConn.Handshake() == nil {
			if err := tlsConn.Close(); err != nil {
				t.Error(err)
			}
		}
	}()

	clientHelloInfo := <-clientHelloChan
	if clientHelloInfo == nil {
		t.Fatal("expected clientHelloInfo, but got nil")
	}
	if clientHelloInfo.ServerName != "example.com" {
		t.Errorf("expected %q, but got %q", "example.com", clientHelloInfo.ServerName)
	}
}
