// mongodb客户端库mgo封装
package mgo

import (
	"gopkg.in/mgo.v2"

	"util/logs"
)

//
type Collection struct {
	*mgo.Session
	*mgo.Collection
}

//func (this *Collection) Close() {
//	this.Session.Close()
//}

//
var g_mdb *mgo.Session

// url:#[mongodb://][user:pass@]host1[:port1][,host2[:port2],…][/database][?options]
func Init(addrs string) {
	session, e := mgo.Dial(addrs)
	if e != nil {
		logs.Panicln(e)
	}
	g_mdb = session
}

//
func GetCollection(dbName, collectionName string) *Collection {
	s := g_mdb.New()
	c := s.DB(dbName).C(collectionName)

	return &Collection{
		Session:    s,
		Collection: c,
	}
}
