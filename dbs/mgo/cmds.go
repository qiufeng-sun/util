/**
Create by SunXiguang 2018-12-28
Desc: mongodb命令定义
*/
package mgo

import (
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

//https://blog.csdn.net/yisun123456/article/details/78591255
//https://docs.mongodb.com/manual/reference/operator/query/

//
const (
	SET     = "$set"
	UNSET   = "$unset"
	INC     = "$inc"
	ADD2SET = "$addToSet"
	EACH    = "$each"
)

////////////////////////////////////////////////////////////////////////////
//
type M struct {
	bson.M
}

//
func NewM() *M {
	return &M{
		M: bson.M{},
	}
}

//
func (this *M) Op(op, key string, val interface{}) *M {
	// check op
	v, ok := this.M[op]
	if !ok {
		v = bson.M{}
		this.M[op] = v
	}

	// check key
	m := v.(bson.M)
	if _, ok := m[key]; ok {
		panic("already set!")
	}

	m[key] = val

	return this
}

//
func (this *M) Set(key string, val interface{}) *M {
	if reflect.ValueOf(val).IsNil() {
		return this
	}
	return this.Op(SET, key, val)
}

//
func (this *M) Unset(key string, val interface{}) *M {
	return this.Op(UNSET, key, val)
}

//
func (this *M) Inc(key string, val interface{}) *M {
	return this.Op(INC, key, val)
}

//
func (this *M) Add2Set(key string, vals ...interface{}) *M {
	num := len(vals)
	if 0 == num {
		panic("Add2Set|no val to set! key:" + key)
	}

	if 1 == num {
		return this.Op(ADD2SET, key, vals[0])
	}

	return this.Op(ADD2SET, key, bson.M{EACH: vals})
}
