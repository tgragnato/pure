package sni

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

func TestCopyLoop(t *testing.T) {
	t.Parallel()

	client, backend := net.Pipe()
	go copyLoop(client, client, backend)
	buf := new(bytes.Buffer)
	go func() {
		_, err := io.Copy(buf, backend)
		if err != nil {
			t.Error(err)
		}
	}()
	testData := []byte("Hello, World!")
	if _, err := client.Write(testData); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	client.Close()
	backend.Close()

	receivedData := buf.Bytes()
	if !bytes.Equal(receivedData, testData) {
		t.Errorf("Expected %q, but got %q", testData, receivedData)
	}
}
