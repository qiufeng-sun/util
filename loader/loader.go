// 读取本地配置文件. 格式:xml,json
package loader

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"

	"util/logs"
)

// 解析函数
type funcParse func(data []byte, v interface{}) error

// 读取并解析文件数据
func ParseFile(parse funcParse, fileName string, out interface{}) error {
	logs.Debug("load file<%s>", fileName)

	// 打开并读取文件
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		logs.Error("file %v load failed! err=%v", fileName, err.Error())
		return err
	}

	// 解析数据
	if err = parse(data, out); err != nil {
		logs.Error("file %v data parse failed! err=%v", fileName, err.Error())
		return err
	}

	return nil
}

// json文件
func ParseJsonFile(fileName string, out interface{}) error {
	return ParseFile(json.Unmarshal, fileName, out)
}

// xml文件
func ParseXmlFile(fileName string, out interface{}) error {
	if err := ParseFile(xml.Unmarshal, fileName, out); err != nil {
		return err
	}

	if entrys, ok := out.(IXmlEntrys); ok {
		entrys.InitMap()
	}

	return nil
}

// xml资源接口
type IXmlEntrys interface {
	InitMap()
}
