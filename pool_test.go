package conn_pool

import (
	"testing"
	"os"
	"net"
	"time"
	"sync"
)

var run = true
var ln net.Listener
var target = []byte("hello world")
var wg = &sync.WaitGroup{}

func TestMain(m *testing.M) {
	createTcpServer()
	retCode := m.Run()
	wg.Wait()
	os.Exit(retCode)
}

func createTcpServer() {
	var err error
	ln, err = net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			if conn, err := ln.Accept(); err != nil {
				panic(err)
			} else {
				if err = conn.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
					panic(err)
				}
				if _, err := conn.Write(target); err != nil {
					panic(err)
				}
			}
		}
	}()
}

func TestPool(t *testing.T) {
	addr := ln.Addr()

	pool, err := NewPool(addr.Network(), addr.String(), 10)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := pool.Get()
	if err != nil {
		t.Fatal(err)
	}

	b := make([]byte, len(target))
	n, err := conn.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(target) {
		t.Fatal(err)
	}

	if err = pool.Put(conn); err != nil {
		t.Fatal(err)
	}
}
