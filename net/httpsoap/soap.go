package httpsoap

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"

	"util/logs"
)

type Envelope struct {
	XMLName       xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Xsi           string   `xml:"xmlns:xsi,attr"`
	Soapenc       string   `xml:"xmlns:soapenc,attr"`
	Xsd           string   `xml:"xmlns:xsd,attr"`
	EncodingStyle string   `xml:"soap:encodingStyle,attr"`
	Soap          string   `xml:"xmlns:soap,attr"`
	Body          Body
}

type Body struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Data    string   `xml:",innerxml"`
}

func newEnvelope(data interface{}) Envelope {
	msg, err := xml.Marshal(data)
	if err != nil {
		panic(err)
	}
	return Envelope{
		Xsi:           "http://www.w3.org/2001/XMLSchema-instance",
		Soapenc:       "http://schemas.xmlsoap.org/soap/encoding/",
		Xsd:           "http://www.w3.org/2001/XMLSchema",
		EncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		Soap:          "http://schemas.xmlsoap.org/soap/envelope/",
		Body:          Body{Data: string(msg)},
	}
}

func writeEnvelope(env Envelope, writer io.Writer) error {
	msg, err := xml.Marshal(env)
	if err != nil {
		return err
	}

	writer.Write(msg)
	return nil
}

func SendEnvelope(data interface{}, url string, action string) (*http.Response, error) {
	//
	buf := new(bytes.Buffer)
	buf.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)

	//
	env := newEnvelope(data)
	err := writeEnvelope(env, buf)
	if err != nil {
		return nil, err
	}
	logs.Debug("envelop:%v\n", buf.String())

	//
	r, err := http.Post(url, "application/soap+xml; action="+action, buf)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func ReadEnvelope(reply interface{}, reader io.Reader) error {
	var b [10000]byte
	n, _ := reader.Read(b[:])
	logs.Debug("return body=%v\n", string(b[:n]))

	var env Envelope
	e := xml.Unmarshal(b[:n], &env)
	logs.Debug("e=%v, env=%v\n", e, env)
	if e != nil {
		return e
	}

	logs.Debug("return=%v\n", env.Body.Data)

	return xml.Unmarshal([]byte(env.Body.Data), reply)
}
