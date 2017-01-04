package etcd

import (
	"encoding/json"
	"testing"

	"fmt"
	"time"
)

//
func toString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

//
func TestClient(t *testing.T) {
	//
	etcdAddrs := []string{"http://127.0.0.1:2379"}

	//
	servers := [][]string{
		{"test1", "127.0.0.1:9901"},
		{"test11", "127.0.0.1:9911"},
		{"test11", "127.0.0.1:9912"},
	}

	//
	value := fmt.Sprintf(`{"reg_time":%v}`, time.Now().Unix())

	// register servers
	var sr = make([]*EtcdRegister, len(servers))

	for i, s := range servers {
		//
		r := NewEtcdRegister(etcdAddrs, time.Second, time.Second*3)
		r.Register(s[0], s[1], value)

		//
		sr[i] = r
	}

	// client
	check := map[string][]*ServerInfo{
		"test1": []*ServerInfo{&ServerInfo{"127.0.0.1:9901", value}},
		"test11": []*ServerInfo{
			&ServerInfo{"127.0.0.1:99011", value},
			&ServerInfo{"127.0.0.1:99012", value},
		},
	}

	for k, v := range check {
		c := NewEtcdClient(etcdAddrs, k)
		infos := c.GetServerInfos()
		t.Logf("name:%v\naddrs:%#v\nexpected addrs:%v\n\n",
			k, toString(infos), toString(v))
	}
}

//
func TestWatch(t *testing.T) {
	//
	etcdAddrs := []string{"http://127.0.0.1:2379"}
	d := []string{"test_watch", "127.0.0.1:1000"}

	//
	c := NewEtcdClient(etcdAddrs, d[0])
	t.Log("servers:\n", toString(c.GetServerInfos()))

	// reg
	r := NewEtcdRegister(etcdAddrs, time.Second, time.Second*3)
	r.debug = true
	r.Register(d[0], d[1], "watching")

	time.Sleep(time.Second * 5)
}

//
func TestWatch2(t *testing.T) {
	//
	etcdAddrs := []string{"http://127.0.0.1:2379"}
	d := []string{"test_watch", "127.0.0.1:1000"}

	//
	c := NewEtcdClient(etcdAddrs, d[0])
	t.Log("servers:\n", toString(c.GetServerInfos()))

	// reg
	r := NewEtcdRegister(etcdAddrs, time.Second, time.Second*3)
	r.debug = true
	r.Register(d[0], d[1], "watching")

	time.Sleep(time.Second)
	r.Unregister()
	time.Sleep(time.Second)
	r.Register(d[0], d[1], "watching")
	time.Sleep(time.Second)
}

//
func TestAddWatchs(t *testing.T) {
	etcdAddrs := []string{"http://127.0.0.1:2379"}
	svc := "test_addwatchs"
	servers := [][]string{
		{svc, "127.0.0.1:9901"},
		{svc, "127.0.0.1:9911"},
		{svc, "127.0.0.1:9912"},
	}

	fUpdate := func(svc string, addrs []string) {
		t.Logf("svc:%v, addrs:%v", svc, addrs)
	}

	AddWatchs(etcdAddrs, []string{svc}, fUpdate)

	for _, info := range servers {
		r := NewEtcdRegister(etcdAddrs, time.Second, time.Second*2)
		r.Register(info[0], info[1], "")

		time.Sleep(time.Second)
		time.AfterFunc(time.Second*2, func() {
			r.Unregister()
		})
	}
	time.Sleep(time.Second * 3)
}
