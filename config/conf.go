package config

import (
	"strings"
	"todo_list/model"

	"gopkg.in/ini.v1"
)

var (
	AppMode        string
	HttpIp         string
	HttpPort       string
	RedisAddr      string
	RedisPw        string
	RedisDbName    string
	Db             string
	DbHost         string
	DbPort         string
	DbUser         string
	DbPassWord     string
	DbName         string
	NgnixImagePort string
)

func Init() {
	file, err := ini.Load("./config/config.ini")
	if err != nil {
		println("配置文件加载失败")
	}
	LoadServer(file)
	LoadMysql(file)
	LoadNgnix(file)
	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&&parseTime=true"}, "")
	model.DataBase(path)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpIp = file.Section("service").Key("HttpIp").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadNgnix(file *ini.File) {
	NgnixImagePort = file.Section("ngnix").Key("ImagePort").String()
}
