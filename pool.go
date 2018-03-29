package conn_pool

import (
	"net"
	"errors"
)

type Pool struct {
	network, address string
	connChan         chan net.Conn
}

func NewPool(network, address string, size int) (*Pool, error) {
	connChan := make(chan net.Conn, size)

	for i := 0; i < size; i++ {
		if conn, err := net.Dial(network, address); err != nil {
			for {
				select {
				case conn := <-connChan:
					conn.Close()
				}
			}
			return nil, err
		} else {
			connChan <- conn
		}
	}

	return &Pool{network, address, connChan}, nil
}

func (p *Pool) Get() (net.Conn, error) {
	select {
	case conn := <-p.connChan:
		return conn, nil
	default:
		return nil, errors.New("no connection available")
	}
}

func (p *Pool) Put(conn net.Conn) (error) {
	select {
	case p.connChan <- conn:
		return nil
	default:
		return errors.New("put conn back fail")
	}
}
