//////////////////////////////////////////////////////////////
//// FileLogWriter implements LoggerInterface.
//// It writes messages by lines limit, file size limit, or time frequency.
//type FileLogWriter struct {
//	*log.Logger
//	mw *MuxWriter
//	// The opened file
//	Filename string `json:"filename"`

//	Maxlines          int `json:"maxlines"`
//	maxlines_curlines int

//	// Rotate at size
//	Maxsize         int `json:"maxsize"`
//	maxsize_cursize int

//	// Rotate daily
//	Daily          bool  `json:"daily"`
//	Maxdays        int64 `json:"maxdays"`
//	daily_opendate int

//	Rotate bool `json:"rotate"`

//	startLock sync.Mutex // Only one log can write to the file

//	Level int `json:"level"`
//}

//const (
//	LevelEmergency = iota
//	LevelAlert
//	LevelCritical
//	LevelError
//	LevelWarning
//	LevelNotice
//	LevelInformational
//	LevelDebug
//)
//////////////////////////////////////////////////////////////

package logs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	ghLogs "github.com/astaxie/beego/logs"

	"runtime"
	"util"
)

//
type LogCfg struct {
	Filename string `json:"filename,omitempty"`
	MaxLines int    `json:"maxlines,omitempty"` // Rotate at line
	MaxSize  int    `json:"maxsize,omitempty"`  // Rotate at size
	Daily    bool   `json:"daily,omitempty"`    // Rotate daily
	MaxDays  int64  `json:"maxdays,omitempty"`
	Rotate   bool   `json:"rotate,omitempty"`
	Level    int    `json:"level,omitempty"`
	Perm     string `json:"perm,omitempty"`
}

// logger references the used application logger.
var g_log = ghLogs.GetBeeLogger()
var g_logLv = ghLogs.LevelEmergency

func GetLogger() *ghLogs.BeeLogger {
	return g_log
}

func Init(cfgFile string) *ghLogs.BeeLogger {
	// 打开并读取文件
	data, e := ioutil.ReadFile(cfgFile)
	if e != nil {
		msg := fmt.Sprintf("file %v load failed! error=%v", cfgFile, e)
		panic(msg)
	}

	var m map[string]*LogCfg
	if e := json.Unmarshal(data, &m); e != nil {
		msg := fmt.Sprintf("log config invalid! file:%v, error:%v\n", cfgFile, e)
		panic(msg)
	}

	for k, v := range m {
		delLogger(k)
		SetLogger(k, v)
		if v.Level > g_logLv {
			g_logLv = v.Level
		}
	}

	g_log.SetLevel(g_logLv)

	return g_log
}

func Close() {
	g_log.Close()
}

func Async(msgNum int64) {
	g_log.Async(msgNum)
}

func SetLogger(adapterName string, cfg *LogCfg) {
	jsonCfg := util.ToJsonString(cfg)
	fmt.Printf("set adapter:%v, cfg:%v\n", adapterName, jsonCfg)

	SetLoggerByJson(adapterName, jsonCfg)
}

func delLogger(adapterName string) {
	g_log.DelLogger(adapterName)
}

func SetLoggerByJson(adapterName, jsonCfg string) {
	e := g_log.SetLogger(adapterName, jsonCfg)
	if e != nil {
		fmt.Printf("set logger failed! error=%v", e)
	}
}

func SetLogLv(lv int) {
	if lv > g_logLv {
		g_logLv = lv
		g_log.SetLevel(g_logLv)
	}
}

func Debug(format string, v ...interface{}) {
	if ghLogs.LevelDebug > g_logLv {
		return
	}

	fixformat := buildCaller(format)
	g_log.Debug(fixformat, v...)
}

func Debugln(v ...interface{}) {
	g_log.Debug(strings.Repeat("%v ", len(v)), v...)
}

func Error(format string, v ...interface{}) {
	fixformat := buildCaller(format)
	g_log.Error(fixformat, v...)
}

func Errorln(v ...interface{}) {
	g_log.Error(strings.Repeat("%v ", len(v)), v...)
}

func Info(format string, v ...interface{}) {
	g_log.Info(format, v...)
}

func Infoln(v ...interface{}) {
	g_log.Info(strings.Repeat("%v ", len(v)), v...)
}

func Warn(format string, v ...interface{}) {
	if ghLogs.LevelWarning > g_logLv {
		return
	}

	fixformat := buildCaller(format)
	g_log.Warning(fixformat, v...)
}

func Warnln(v ...interface{}) {
	g_log.Warning(strings.Repeat("%v ", len(v)), v...)
}

func Critical(format string, v ...interface{}) {
	g_log.Critical(format, v...)
}

func Panicln(v ...interface{}) {
	s := fmt.Sprintf(strings.Repeat("%v ", len(v)), v...)
	panic(s)
}

func buildCaller(format string) string {
	_, file, line, _ := runtime.Caller(2)
	callerformat := fmt.Sprintf("[%v:%v] ", file, line) + format
	return callerformat
}
