package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
)

// Copied from golint
var commonInitialisms = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UI", "UID", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
var commonInitialismsReplacer *strings.Replacer

func init() {
	var commonInitialismsForReplacer []string
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

var smap = map[string]string{}

func ToDBName(name string) string {
	if v, ok := smap[name]; ok {
		return v
	}

	value := commonInitialismsReplacer.Replace(name)
	buf := bytes.NewBufferString("")
	for i, v := range value {
		if i > 0 && v >= 'A' && v <= 'Z' {
			buf.WriteRune('_')
		}
		buf.WriteRune(v)
	}

	s := strings.ToLower(buf.String())
	smap[name] = s
	return s
}

func RandStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func Unmarshal(str string, obj interface{}) {
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		panic("Marshalling failed: " + err.Error())
	}
}

func MarshalToStr(obj interface{}) string {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		panic("Marshalling failed: " + err.Error())
	}

	return string(objBytes)
}

func Dump(obj interface{}) {
	fmt.Println("Dump:", obj)
}

func DumpPanic(obj interface{}) {
	panic(fmt.Sprintf("Dump:", obj))
}
