package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
)

////////////////////////////////////////////////////////////////////////////////
//
func Unzip(zipDatas []byte) (string, error) {
	//
	in := bytes.NewReader(zipDatas)
	gzipReader, e := gzip.NewReader(in)
	if e != nil {
		return "", e
	}
	defer gzipReader.Close()

	//
	out := &bytes.Buffer{}

	//
	io.Copy(out, gzipReader)

	return out.String(), nil
}
