package loader

import (
	"testing"
)

// 物品静态属性
type ItemEntry struct {
	//	XMLName     xml.Name 		`xml:"item"`
	Id   int `xml:"id,attr"`
	Type int `xml:"type,attr"`
	Qlty int `xml:"qlty,attr"`
}

// 物品资源管理对象定义
type ItemEntrys struct {
	//	XMLName     xml.Name 		`xml:"root"`
	Entrys []ItemEntry `xml:"item"`
	entrys map[int]*ItemEntry
}

// 资源接口实现
// IEntrys
func (this *ItemEntrys) InitMap() {
	this.entrys = make(map[int]*ItemEntry, len(this.Entrys))
	for _, entry := range this.Entrys {
		this.entrys[entry.Id] = &entry
	}
}

//
func TestXmlLoader(t *testing.T) {
	//
	var entrys ItemEntrys

	//
	e := ParseXmlFile("example.xml", &entrys)
	t.Logf("error:%v\nitem entrys:%#v", e, entrys)

	if e != nil {
		t.Error(e)
	}
}

type JsonCfg struct {
	Typ   int         `json:"type"`
	Desc  string      `json:"desc,omitempty"`
	Title string      `json:"title"`
	Sub   *SubJsonCfg `json:"sub"`
}

type SubJsonCfg struct {
	SubTile string `json:"title,omitempty"`
}

//
func TestJsonLoader(t *testing.T) {
	//
	var cfg JsonCfg

	//
	e := ParseJsonFile("example.json", &cfg)
	t.Logf("error:%v\njson:%#v", e, cfg)

	if e != nil {
		t.Error(e)
	}
}
