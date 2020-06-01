package iniconfig

import (
	"io/ioutil"
	"testing"
)

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
	Port     int    `ini:"port"`
}

func TestIniConfig(t *testing.T) {

	data, err := ioutil.ReadFile("../config.ini")
	if err != nil {
		t.Error("read file failed")
	}
	var conf Config
	err = UnMarshal(data, &conf)
	if err != nil {
		t.Errorf("unmarshalk failed,err:%v", err)
		return
	}

	_, err = Marshal(conf)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("unmarshal success, conf:%#v", conf)
}
