package etcd

import (
	"testing"

	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
)

////////////////////////////////////////////////////////////////////////////////
//
func TestRegister(t *testing.T) {
	//
	etcdAddrs := []string{"http://127.0.0.1:2379"}

	//
	servers := [][]string{
		{"test1", "127.0.0.1:9901"},
		{"test11", "127.0.0.1:9911"},
		{"/test11", "127.0.0.1:9912"},
	}

	//
	value := fmt.Sprintf(`{"reg_time":%v}`, time.Now().Unix())

	// register servers
	for _, s := range servers {
		//
		r := NewEtcdRegister(etcdAddrs, time.Second, time.Second*3)
		r.Register(s[0], s[1], value)
	}

	// output server nodes
	r := NewEtcdRegister(etcdAddrs, time.Second, time.Second*3)
	option := &client.GetOptions{Recursive: true}
	for i := 0; i < 5; i++ {
		resp, e := r.KeysAPI.Get(context.Background(), servers[0][0], option)
		t.Logf("test1 - %v: error=%v, resp=%v\n", i, e, toString(resp))

		resp, e = r.KeysAPI.Get(context.Background(), servers[1][0], option)
		t.Logf("test11 - %v: error=%v, %v\n\n", i, e, toString(resp))

		time.Sleep(time.Second)
	}
}

//
func TestUnregister(t *testing.T) {
	//
	etcdAddrs := []string{"http://127.0.0.1:2379"}
	serverPath := "test_unreg"
	serverAddr := "127.0.0.1:9999"

	//
	r := NewEtcdRegister(etcdAddrs, time.Second, time.Second*3)
	r.Register(serverPath, serverAddr, "")

	option := &client.GetOptions{Recursive: true}
	resp, e := r.KeysAPI.Get(context.Background(), serverPath, option)
	t.Logf("reg: error=%v, %v\n\n", e, toString(resp))

	//
	r.Unregister()

	resp, e = r.KeysAPI.Get(context.Background(), serverPath, option)
	t.Logf("unreg: error=%v, %v\n\n", e, toString(resp))
}

// 测试心跳异常时, 是否可以重新注册
//func TestDisconnect(t *testing.T) {
//	//
//	etcdAddrs := []string{"http://127.0.0.1:2379"}
//	serverPath := "test_disconnect"
//	serverAddr := "127.0.0.1:9999"

//	//
//	r := NewEtcdRegister(etcdAddrs, time.Millisecond*200, time.Second)
//	r.debug = true
//	r.Register(serverPath, serverAddr, "")

//	time.Sleep(time.Second * 500)
//}
