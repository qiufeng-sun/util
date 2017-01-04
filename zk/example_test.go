package zk

import (
	"encoding/json"
	"sync"
	"testing"
	"time"
)

//
type Cfg struct {
	Num int
}

type AtomicCfg struct {
	*sync.Mutex
	*Cfg
}

func (this *AtomicCfg) Unmarshal(data []byte) (interface{}, error) {
	var cfg *Cfg

	if e := json.Unmarshal(data, &cfg); e != nil {
		return nil, e
	}

	this.Lock()
	defer this.Unlock()

	this.Cfg = cfg

	return cfg, nil
}

var g_cfg = &AtomicCfg{Mutex: &sync.Mutex{}}

func GetCfg() *Cfg {
	g_cfg.Lock()
	r := g_cfg.Cfg
	g_cfg.Unlock()

	return r
}

//
func LoadCfg() {
	AddZkParser("/services/game", g_cfg)

	LoadZkDatasW([]string{"localhost:1280"}, time.Second*5)
}

//
func TestCfg(t *testing.T) {}
