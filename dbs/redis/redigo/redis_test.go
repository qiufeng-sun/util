package redigo

import (
	"testing"

	"math/rand"

	ghRedis "github.com/garyburd/redigo/redis"
)

//
func TestHealthCheck(t *testing.T) {
	e := HealthCheck()
	t.Logf("health check: error=%v\n", e)

	if e != nil {
		t.Error(e)
	}
}

//
func TestRedis(t *testing.T) {
	c := getSampleConn("")
	defer c.Close()

	k := "testing:unit"
	v := rand.Int()

	t.Logf("input: k=%v, v=%v\n", k, v)

	r, e := c.Do("SET", k, v)
	t.Logf("set res: k=%v, v=%v, reply=%v, error=%v\n", k, v, r, e)

	if e != nil {
		t.Error(e)
		return
	}

	r1, e := ghRedis.Int(c.Do("GET", k))
	t.Logf("get res: k=%v, v=%v, error=%v\n", k, r1, e)

	if e != nil {
		t.Error(e)
		return
	}

	if r1 != v {
		t.Errorf("test failed! k=%v, res=%v, expected=%v\n", k, r1, v)
		return
	}
}
