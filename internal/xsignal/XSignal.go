package xsignal

import (
	"sync"

	"github.com/quasilyte/gsignal"
)

type Xconnection struct {
	dis bool
}

func (c *Xconnection) IsDisposed() bool {
	return c.dis
}

func (c *Xconnection) Disconnect() {
	c.dis = true
}

type XConnector struct {
	mu    sync.Mutex
	conns map[string]func() // 存储断开函数
}

func NewXConnector() *XConnector {
	return &XConnector{conns: make(map[string]func())}
}

func (c *XConnector) Disconnect(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if fn, ok := c.conns[id]; ok {
		fn() // 立即断开
		delete(c.conns, id)
	}
}

func (c *XConnector) DisconnectAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for id, fn := range c.conns {
		fn()
		delete(c.conns, id)
	}
}

func ConnectSlot[T any](c *XConnector, sig *gsignal.Event[T], slot func(T), id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	conn := &Xconnection{} // 仍可作为标识
	if sig == nil {
		return
	}
	sig.Connect(conn, slot)
	c.conns[id] = func() {
		sig.Disconnect(conn) // 立即断开
	}
}

func ConnectSig[T any](c *XConnector, sigSend *gsignal.Event[T], sigRecv *gsignal.Event[T], id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn := &Xconnection{}
	sigSend.Connect(conn, func(arg T) {
		sigRecv.Emit(arg)
	})

	c.conns[id] = func() {
		sigSend.Disconnect(conn)
	}
}
