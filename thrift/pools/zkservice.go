package pools

import (
	"errors"
	"sync"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/samuel/go-zookeeper/zk"

	"util"
	"util/logs"
	uzk "util/zk"
)

var (
	ErrZkEnv       = errors.New("zookeeper environment error")
	ErrZkConnect   = errors.New("zookeeper connection error")
	ErrZkNotExists = errors.New("zookeeper data not exists")
)

//
type ThriftPool struct {
	Server  string
	Percent int
	Weight  float64
}

type ThriftPools struct {
	env        string
	name       string
	level      int
	updateTime int64

	Pools []ThriftPool
	mu    sync.Mutex

	zkCon *zk.Conn
	zkEvt <-chan zk.Event
}

//
func NewThriftPools(env, serviceName string, serviceLevel int) *ThriftPools {
	return &ThriftPools{
		env:        env,
		name:       serviceName,
		level:      serviceLevel,
		updateTime: util.CurMillisecond(),
	}
}

//
func (this *ThriftPools) Init() {
	//
	server, ok := uzk.GetZkServers(this.env)
	if !ok {
		logs.Panicln("zookeeper env wrong:", this.env)
	}

	//
	c, evt, err := zk.Connect(server, time.Second*3)
	if err != nil {
		logs.Panicln(ErrZkConnect, err)
	}

	//
	this.zkCon = c
	this.zkEvt = evt

	//
	if err := this.loadNodes(); err != nil {
		logs.Warnln("zookeeper load thrift pool failed! err: ", err)
	}

	//
	go this.onZkEvent()
}

//
func (this *ThriftPools) loadNodes() error {
	//
	c := this.zkCon

	//
	servicePool := this.name + "/Pool"
	if ok, _, _ := c.Exists(servicePool); !ok {
		return ErrZkNotExists
	}

	list, _, _, err := c.ChildrenW(servicePool)
	if err != nil || len(list) <= 0 {
		return ErrZkNotExists
	}

	//
	pools, err := getWeightPool(c, list, servicePool, this.level)
	if err != nil {
		return err
	}

	total := 0.0
	for _, v := range pools {
		total += v.Weight
	}

	// 权重处理
	if total > 1.0 {
		percent := 0.0
		for i := 0; i < len(pools); i++ {
			p := &pools[i]
			percent += p.Weight
			p.Percent = (int)(100 * percent / total)
		}
	}

	//
	this.mu.Lock()
	defer this.mu.Unlock()
	this.Pools = pools
	this.updateTime = util.CurMillisecond()

	return nil
}

//
func (this *ThriftPools) onZkEvent() {
	for evt := range this.zkEvt {
		if evt.Err != nil {
			continue
		}

		logs.Info("zk event: path=%v, type=%v, state=%v\n", evt.Path, evt.Type, evt.State)

		switch evt.Type {
		case zk.EventSession:
			if zk.StateHasSession == evt.State {
				this.loadNodes()
			}

		case zk.EventNodeChildrenChanged:
			this.loadNodes()
		}
	}
}

func (this *ThriftPools) GetNodes() []ThriftPool {
	this.mu.Lock()
	defer this.mu.Unlock()

	return this.Pools
}

func (this *ThriftPools) GetUpdateTime() int64 {
	return this.updateTime
}

func getWeightPool(c *zk.Conn, list []string, servicePool string, serviceLevel int) ([]ThriftPool, error) {
	result := make([]ThriftPool, 0, len(list))
	for _, pool := range list {
		path := servicePool + "/" + pool
		content, _, _ := c.Get(path)

		iniconf, err := config.NewConfigData("ini", content)
		if err != nil {
			logs.Warn("Failed to parse pool data: %s", err)
			return nil, err
		}
		level, _ := iniconf.Int("server.service.level")
		if serviceLevel == level {
			weight, _ := iniconf.Float("weight")
			t := ThriftPool{Server: pool, Percent: 0, Weight: weight}
			result = append(result, t)
		}

		logs.Infoln(time.Now(), "service pool in zk,", servicePool, level, pool)
	}

	return result, nil
}
