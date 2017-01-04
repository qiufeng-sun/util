package scribe

import (
	"fmt"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/astaxie/beego/config"

	"util"
	"util/logs"
	"util/thrift/scribe"
)

var (
	g_transportFactory = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	g_protocolFactory  = thrift.NewTBinaryProtocolFactoryDefault()
)

type ScribeClient struct {
	hostPort string
	thrift.TTransport
	sc            *scribe.ScribeClient
	msgs          <-chan *scribe.LogEntry
	lastUnsentMsg *scribe.LogEntry
	quit          chan bool // 阻塞
	quitLoop      chan bool // 非阻塞
	index         int       // for test
}

func (this *ScribeClient) Close() {
	// to stop loop
	this.quitLoop <- true

	// wait loop
	<-this.quit

	// close channel
	close(this.quit)
	close(this.quitLoop)
}

func (this *ScribeClient) loop() {
	for {
		if this.connect() {
			this.sendMsgs()
			this.close()
		} else {
			time.Sleep(time.Second)
		}

		select {
		case <-this.quitLoop:
			logs.Debug("receve quit loop")
			this.quit <- true
			return
		default:
		}
	}
}

func (this *ScribeClient) connect() bool {
	socket, err := thrift.NewTSocket(this.hostPort)
	if err != nil {
		logs.Warn("%v: Failed when call NewTSocket()! err=%v\n", util.Caller(0), err.Error())
		return false
	}

	transport := g_transportFactory.GetTransport(socket)
	if err = transport.Open(); err != nil {
		transport.Close()
		logs.Warn("%v: Failed when call GetTransport()! err=%v\n", util.Caller(0), err.Error())
		return false
	}

	this.TTransport = transport
	this.sc = scribe.NewScribeClientFactory(transport, g_protocolFactory)

	return true
}

func (this *ScribeClient) sendMsgs() {
	//
	if this.lastUnsentMsg != nil {
		if _, e := this.sc.Log([]*scribe.LogEntry{this.lastUnsentMsg}); e != nil {
			logs.Debug("send the 1ast unsend msg failed!\n")
			return
		}
		this.lastUnsentMsg = nil
	}

	//
	for {
		select {
		case msg, ok := <-this.msgs:
			if !ok {
				logs.Info("all msgs have been sent\n")
				this.quitLoop <- true
				return
			}

			if g_debug {
				logs.Info("index:%v\ncategory:%v\nmsg:%v",
					this.index, msg.Category, msg.Message)
			}

			_, e := this.sc.Log([]*scribe.LogEntry{msg})
			if e != nil {
				this.lastUnsentMsg = msg

				logs.Warn("scribe log send failed! err=%v", e)
				return
			}
		}
	}
}

func (this *ScribeClient) close() {
	if this.TTransport != nil {
		this.TTransport.Close()
	}
}

func NewScribeClient(hostPort string, chMsg <-chan *scribe.LogEntry, index int) *ScribeClient {
	client := &ScribeClient{
		hostPort:   hostPort,
		TTransport: nil,
		sc:         nil,
		msgs:       chMsg,
		quit:       make(chan bool),
		quitLoop:   make(chan bool, 2),
		index:      index,
	}

	go func() {
		client.loop()
	}()

	return client
}

//
var (
	x_scribeNum = 5
	x_msgBuffSz = 10000

	g_scribeClient []*ScribeClient
	g_chMsgs       chan *scribe.LogEntry

	g_debug = false // for test
)

func Init(logMsg string, confd config.Configer) {
	logs.Info("scribe<%v> init start!", logMsg)
	defer logs.Info("scribe<%v> init end!", logMsg)

	open := confd.DefaultBool("scribe_open", false)
	logs.Info("scribe<%v> open: %v", logMsg, open)
	if !open {
		return
	}
	addr := confd.DefaultString("scribe_addr", "localhost:7915")
	goNum := confd.DefaultInt("scribe_gonum", -1)
	bufNum := confd.DefaultInt("scribe_bufnum", -1)

	InitScribe(addr, goNum, bufNum)
}

func InitScribe(hostPort string, scribeNum, msgBuffSz int) {
	//
	if scribeNum > 0 && scribeNum < 10 {
		x_scribeNum = scribeNum
	}

	//
	if msgBuffSz > 0 && msgBuffSz < 100000 {
		x_msgBuffSz = msgBuffSz
	}

	//
	g_chMsgs = make(chan *scribe.LogEntry, x_msgBuffSz)

	//
	g_scribeClient = make([]*ScribeClient, x_scribeNum)

	//
	for i := 0; i < x_scribeNum; i++ {
		g_scribeClient[i] = NewScribeClient(hostPort, g_chMsgs, i)
	}

	logs.Info("scribe init ok! host: %v, gonum: %v, bufnum: %v",
		hostPort, x_scribeNum, x_msgBuffSz)
}

func CloseScribe() {
	close(g_chMsgs)

	for i, sc := range g_scribeClient {
		logs.Debug("scribe %v closing...", i)
		sc.Close()
		logs.Debug("scribe stop ok")
	}
}

//
func Log(category, format string, v ...interface{}) {
	//
	msg := format
	if len(v) > 0 {
		msg = fmt.Sprintf(format, v...)
	}

	//
	if len(g_scribeClient) <= 0 {
		logs.Info("category:%v\nmsg:%v", category, msg)
		return
	}

	//
	entry := &scribe.LogEntry{
		Category: category,
		Message:  msg,
	}

	select {
	case g_chMsgs <- entry:
	default:
		logs.Warn("scribe log buff is full! discard -- %v:\n%v", category, msg)
	}
}
