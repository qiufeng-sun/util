/**
Create by SunXigaung 2019-01-25
Desc: mgo常用方法封装
*/
package mgo

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"util"
	"util/logs"
)

//
type idGenRes struct {
	CurMaxId int `bson:"curMaxId"`
}

//
func GenId(c *Collection, moduleName string, num int) (int, error) {
	//
	sel := bson.M{"module": moduleName}
	change := mgo.Change{
		Update:    NewM().Inc("curMaxId", num).M,
		Upsert:    true,
		ReturnNew: true,
	}

	var res idGenRes
	_, e := c.Find(sel).Apply(change, &res)
	if e != nil {
		return 0, fmt.Errorf("GenId|db op failed! module:%v, change:%v, error:%v",
			moduleName, util.ToJsonString(change), e)
	}

	logs.Info("GenId|module:%v, curMaxId:%v", moduleName, res.CurMaxId)

	return res.CurMaxId, nil
}
