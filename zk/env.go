// zk配置，需要外部进行初始化
package zk

//
var (
	g_zkEnvKey    = map[string]string{}
	g_zkHostPorts = map[string][]string{}
)

//
func SetZkEnvKey(env, key string) {
	g_zkEnvKey[env] = key
}

func GetZkEnvKey(env string) (string, bool) {
	key, ok := g_zkEnvKey[env]
	return key, ok
}

//
func SetZkServers(key string, servers []string) {
	g_zkHostPorts[key] = servers
}

func GetZkHostPorts(key string) []string {
	return g_zkHostPorts[key]
}

func GetZkServers(key string) ([]string, bool) {
	r, ok := g_zkHostPorts[key]
	return r, ok
}

//var (
//	g_zkEnvKey = map[string]string{
//		"staging":       "staging",
//		"production-sd": "shangdi",
//	}
//)

//func init() {
//	g_zkHostPorts = make(map[string][]string)

//	g_zkHostPorts["staging"] = []string{
//		"zk1.staging.srv:2181",
//		"zk2.staging.srv:2181",
//		"zk3.staging.srv:2181",
//		"zk4.staging.srv:2181",
//		"zk5.staging.srv:2181",
//	}
//}
