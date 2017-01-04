package pools

import (
	"errors"

	"util/logs"
)

type Client interface {
	Close() error
	GetCreateTime() int64
}

// factory method to create new connections
type Factory func() (Client, error)

type Pool interface {
	Get() (Client, error)
	GetByTime(time int64) (Client, error)
	Close()
	Put(Client) error
	Len() int
}

var (
	ErrClosed       = errors.New("pool is closed")
	ErrCapacity     = errors.New("invalid capacity settings")
	ErrClientIsNull = errors.New("connection is nil. rejecting")
)

var debug = false

func SetDebug(d bool) {
	debug = d
}

func NewPool(initialCap, maxCap int, factory Factory) (Pool, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, ErrCapacity
	}

	pool := &channelPool{
		conns:   make(chan Client, maxCap),
		factory: factory,
	}

	for i := 0; i < initialCap; i++ {
		if client, err := factory(); err == nil {
			pool.conns <- client
		}
	}
	if debug {
		logs.Infoln("New channelPool", pool)
	}
	return pool, nil
}
