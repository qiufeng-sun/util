package tempcache

import (
	"testing"

	"fmt"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
//
func init() {
	g_debug = true
}

//
type CacheTest struct {
}

func (this *CacheTest) Load(keys []string) ([]interface{}, error) {
	ret := make([]interface{}, len(keys))
	for i, v := range keys {
		ret[i] = fmt.Sprintf("key_%v_value_%v", v, i+1)
	}
	return ret, nil
}

func (this *CacheTest) GetExpiredSec() int64 {
	return 2
}

func (this *CacheTest) LogName() string {
	return "testing cache"
}

//
func TestGetCache(t *testing.T) {
	ckey := "key_testing"
	vkeys := [][]string{
		[]string{"1", "2", "3", "5"},
		[]string{"21", "5", "23", "25"},
		[]string{"1", "3", "5", "23"},
	}

	Register(ckey, &CacheTest{})

	c := GetCache(ckey)
	for i := 0; i < len(vkeys); i++ {
		k := vkeys[i]
		values := c.GetValues(k)
		t.Logf("%v\n", values)

		time.Sleep(time.Second)
	}
}
