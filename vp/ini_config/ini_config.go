package iniconfig

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Marshal(v interface{}) (data []byte, err error) {

	return
}

/*
[server]
ip = 172.16.0.15
port = 8080

[mysql]
username = root
password = root
database = test
host = 172.16.0.1
port = 3306


type Config struct {
	ServerConf Server `ini:"server"`
	MysqlConf  Mysql  `ini:"mysql"`
}

type Server struct {
	IP   string `ini:"ip"`
	Port string `ini:"port"`
}

type Mysql struct {
	UserName string `ini:"username"`
	Password string `ini:"password"`
	Database string `ini:"database"`
	Host     string `ini:"host"`
	Port     string `ini:"port"`
}
*/

func UnMarshal(data []byte, v interface{}) (err error) {

	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if rt.Kind() != reflect.Ptr || rv.IsNil() {
		err = errors.New("please pass valid pointer")
		return
	}
	if rv.Elem().Kind() != reflect.Struct {
		err = errors.New("please pass struct")
		return
	}

	lineArr := strings.Split(string(data), "\n")

	var lastSectionName string
	for i, line := range lineArr {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line[0] == ';' || line[0] == '#' {
			continue
		}

		if line[0] == '[' {
			if len(line) <= 2 || line[len(line)-1] != ']' {
				err = fmt.Errorf("syntax error, invalid section:%s, lineNo:%d\n", line, i+1)
				return
			}
		}

		if line[0] == '[' && line[len(line)-i] == ']' {
			sectionName := strings.TrimSpace(line[1 : len(line)-1])
			if len(sectionName) == 0 {
				err = fmt.Errorf("syntax error, invalid section:%s, lineNo:%d\n", line, i+1)
				return
			}

			for i := 0; i < rt.Elem().NumField(); i++ {
				filed := rt.Elem().Field(i)
				tagVal := filed.Tag.Get("ini")
				if tagVal == sectionName {
					lastSectionName = filed.Name
					fmt.Println(lastSectionName)
					break
				}
			}

		}
	}

	return
}
