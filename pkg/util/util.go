package util

import (
	"bytes"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	TimeZone   = "Asia/Shanghai"
)

func Str2Camel(text []byte) string {
	arrs := bytes.Split(bytes.ToLower(text), []byte("_"))
	if len(arrs) < 2 {
		return string(bytes.ToLower(text))
	}
	var buf bytes.Buffer
	buf.Write(arrs[0])
	for _, elem := range arrs[1:] {
		buf.Write(bytes.ToUpper([]byte{elem[0]}))
		buf.Write(elem[1:])
	}
	return buf.String()
}
