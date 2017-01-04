package pools

import (
	"util/logs"
)

type channelPool struct {
	conns   chan Client
	factory Factory
}

func (pool *channelPool) Get() (Client, error) {
	conns := pool.conns

	select {
	case client := <-conns:
		if debug {
			logs.Infoln("Get Reuse", client)
		}

		return client, nil

	default:
	}

	client, e := pool.factory()
	if e != nil {
		return nil, e
	}
	if debug {
		logs.Infoln("Get New", client)
	}
	return client, nil
}

func (pool *channelPool) GetByTime(updateTime int64) (Client, error) {
	conns := pool.conns

	select {
	case client := <-conns:
		createTime := client.GetCreateTime()
		if createTime < updateTime {
			if debug {
				logs.Infoln("Get Reuse", client, "expired, close it!")
			}
			client.Close()
		} else {
			if debug {
				logs.Infoln("Get Reuse", client)
			}
			return client, nil
		}

	default:
	}

	client, e := pool.factory()
	if e != nil {
		return nil, e
	}
	if debug {
		logs.Infoln("Get New", client)
	}
	return client, nil
}

func (pool *channelPool) Put(client Client) error {
	if client == nil {
		return ErrClientIsNull
	}

	select {
	case pool.conns <- client:
		return nil
	default:
		// pool is full
		return client.Close()
	}
}

func (pool *channelPool) Close() {
	conns := pool.conns

	close(conns)
	for client := range conns {
		client.Close()
	}
}

func (pool *channelPool) Len() int {
	return len(pool.conns)
}
