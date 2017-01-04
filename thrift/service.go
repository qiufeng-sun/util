package thrift

import (
	"fmt"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"

	"util"
	"util/logs"
)

var Debug bool = false

type ThriftService struct {
	Transport  thrift.TTransport
	Server     string
	isClose    bool
	CreateTime int64
}

func (s *ThriftService) GetCreateTime() int64 {
	return s.CreateTime
}

func (s *ThriftService) Close() error {
	s.isClose = true
	s.Transport.Close()
	return nil
}

func (s *ThriftService) IsClose() bool {
	return s.isClose
}

func (s *ThriftService) Check(method func() (interface{}, error)) (interface{}, error) {
	// check run time
	defer logs.CheckTime(500, 1, Debug)()

	result, err := method()
	if err != nil {
		if _, ok := err.(thrift.TProtocolException); ok {
			// TProtocolException 发生thrift协议层错误,关闭连接
			s.Close()
		}
	}

	return result, err
}

func NewThriftService(server string) (*ThriftService, error) {
	return NewThriftServiceWithTimeout(server, time.Second*5)
}

func NewThriftServiceWithTimeout(server string, timeout time.Duration) (*ThriftService, error) {
	socket, err := thrift.NewTSocket(server)
	socket.SetTimeout(timeout)
	if err != nil {
		return nil, fmt.Errorf("thrift error resolving address: %s, %s", server, err)
	}

	transport := thrift.NewTFramedTransport(socket)
	if err := transport.Open(); err != nil {
		logs.Warnln("connect error:", server, err)
		return nil, fmt.Errorf("thrift error connect server: %s", server, err)
	}

	time := util.CurMillisecond()
	return &ThriftService{Transport: transport, Server: server, CreateTime: time}, nil
}

func NewSimpleServer(processor thrift.TProcessor, addr string) *thrift.TSimpleServer {
	socket, err := thrift.NewTServerSocket(addr)
	if err != nil {
		logs.Warnln(err)
		return nil
	}

	transportFactory := thrift.NewTTransportFactory()
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	server := thrift.NewTSimpleServer4(processor, socket, transportFactory, protocolFactory)
	return server
}
