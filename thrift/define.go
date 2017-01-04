package thrift

//
type DeployEnv map[string]string

func NewDeployEnv() DeployEnv {
	return make(map[string]string)
}

func (this DeployEnv) Set(env, realEnv string) {
	this[env] = realEnv
}

func (this DeployEnv) RealEnv(env string) string {
	r, ok := this[env]
	if ok {
		return r
	}
	return this["default"]
}
