package database

import (
	"fmt"

	"github.com/astaxie/beego/orm"

	"core/time"
	"util/loader"
	"util/logs"
)

//
var (
	g_dbCfgs = make(map[string]DatabaseConfig) // <aliasName, cfg>
)

//
type DatabaseConfig struct {
	Driver       string `json:"dbdriver"`
	Protocol     string `json:"protocol"`
	Address      string `json:"address"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Database     string `json:"database"`
	Charset      string `json:"charset"`
	MaxIdle      int    `json:"maxIdle"`
	MaxConn      int    `json:"maxConn"`
	CheckSec     int    `json:"checkSec"`
	CurCheckTime int64
}

func (this DatabaseConfig) GenDsn() string {
	switch this.Driver {
	case "mysql":
		username := this.Username
		if this.Password != "" {
			username = username + ":" + this.Password
		}
		dsn := fmt.Sprintf("%s@%s(%s)/%s?charset=%s",
			username,
			this.Protocol,
			this.Address,
			this.Database,
			this.Charset)

		return dsn + "&loc=Asia%2FShanghai"
	}

	panic(fmt.Sprintf("unsupport db driver. driver=%v\n", this.Driver))
}

//
func InitByFile(fileName string) {
	cfgs := make(map[string]DatabaseConfig)
	if err := loader.ParseJsonFile(fileName, &cfgs); err != nil {
		panic(fmt.Sprintf("Error:%v\n", err.Error()))
	}

	InitByCfgs(cfgs)
}

//
func InitByCfgs(cfgs map[string]DatabaseConfig) {
	for k, cfg := range cfgs {
		if _, ok := g_dbCfgs[k]; ok {
			panic(fmt.Sprintf("duplicate alias name! duplicate:%v\n", k))
		}

		cfg.CurCheckTime = time.Now().Unix()
		g_dbCfgs[k] = cfg

		dsn := cfg.GenDsn()
		orm.RegisterDataBase(k, cfg.Driver, dsn, cfg.MaxIdle, cfg.MaxConn)

		logs.Debug("db init! alias:%v, dsn:%v\n", k, dsn)
	}
}

//
const x_databaseCfgFile = "conf/database.conf"

func Init() {
	InitByFile(x_databaseCfgFile)
}

//
func GetOrm(name string) orm.Ormer {
	o := orm.NewOrm()
	o.Using(name)

	return o
}

//
func GetDefOrm() orm.Ormer {
	return orm.NewOrm()
}

//
const sql_ping = "select 1"

//
func HealthCheck() error {

	nowUnix := time.Now().Unix()

	o := orm.NewOrm()
	for k, dbCfg := range g_dbCfgs {
		if nowUnix-dbCfg.CurCheckTime >= int64(dbCfg.CheckSec) {

			logs.Debug("HealthCheck| Test Mysql name=%v", k)

			o.Using(k)
			_, e := o.Raw(sql_ping).Exec()
			if e != nil {
				logs.Error("HealthCheck|test mysql failed! name=%v, err=%v", k, e)
			}

			dbCfg.CurCheckTime = nowUnix
			g_dbCfgs[k] = dbCfg
		}
	}

	return nil
}
