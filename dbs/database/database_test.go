package database

import (
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//
func TestInsert(t *testing.T) {
	o := orm.NewOrm()

	sql := "insert into test_go_orm set id=?, test=?"
	//sql := "insert into test_go_orm(id,test) values(?,?)"

	res, err := o.Raw(sql, 1, "name2").Exec()

	t.Log("sql=", sql)
	rows, _ := res.RowsAffected()
	t.Log("err=", err, "row effected=", rows)
}

//
func TestSelect(t *testing.T) {
	o := orm.NewOrm()

	sql := "select * from test_go_orm where `id`=?"

	var maps []orm.Params
	num, err := o.Raw(sql, 1).Values(&maps)

	t.Log("sql=", sql)
	t.Log("num=", num, "err=", err)
	if err == nil && num > 0 {
		t.Log("select ok! id=%v, test=%v\n", maps[0]["id"], maps[0]["test"])
	} else {
		t.Error("select failed!\n")
	}
}

//
func TestHealthCheck(t *testing.T) {
	if e := HealthCheck(); e != nil {
		t.Fatal("TestHealthCheck(): failed!", e)
	} else {
		t.Fatal("TestHealthCheck(): ok")
	}
}

//
func init() {
	InitByFile("mysql.conf")
}
