package sni

import (
	"context"
	"io"
	"net"
	"sync"
)

type readerCtx struct {
	ctx context.Context
	r   io.Reader
}

func (r *readerCtx) Read(p []byte) (n int, err error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}

func safeCopy(dst net.Conn, src io.Reader, wg *sync.WaitGroup, ctx context.Context, cancel context.CancelFunc) {
	defer wg.Done()
	r := &readerCtx{ctx: ctx, r: src}
	_, err := io.Copy(dst, r)
	dst.(*net.TCPConn).CloseWrite()
	if err != nil {
		cancel()
	}
}

func copyLoop(clientR io.Reader, clientW net.Conn, backend net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())
	go safeCopy(clientW, backend, &wg, ctx, cancel)
	go safeCopy(backend, clientR, &wg, ctx, cancel)
	wg.Wait()
}
