package pools

import "math/rand"

type ZooData struct {
	*ThriftPools
}

func NewZkData(env, name string, level int) *ZooData {
	data := &ZooData{
		ThriftPools: NewThriftPools(env, name, level),
	}

	data.ThriftPools.Init()

	return data
}

func (z *ZooData) GetThriftNode() string {
	d := z.GetNodes()
	if len(d) <= 0 {
		return ""
	}
	per := rand.Intn(100)
	for _, v := range d {
		if per <= v.Percent {
			return v.Server
		}
	}

	return d[0].Server
}
