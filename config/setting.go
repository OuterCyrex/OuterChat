package config

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
)

var (
	JwtKey string
)

var (
	Addr        string
	DB          int
	PoolSize    int
	MinIdleConn int
)

func LoadDatabase(file *ini.File) {
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassword = file.Section("database").Key("DbPassword").MustString("***")
	DbHost = file.Section("database").Key("DbHost").MustString("127.0.0.1")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbName = file.Section("database").Key("DbName").MustString("outerchat")
}

func LoadJwtKey(file *ini.File) {
	JwtKey = file.Section("jwt").Key("jwtKey").MustString("12q3ewe565y54535")
}

func LoadCache(file *ini.File) {
	Addr = file.Section("cache").Key("Addr").MustString("127.0.0.1:6379")
	DB = file.Section("cache").Key("DB").MustInt(0)
	PoolSize = file.Section("cache").Key("PoolSize").MustInt(30)
	MinIdleConn = file.Section("cache").Key("MinIdleConn").MustInt(30)
}

func InitConfig() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("初始化出错", err.Error())
	}
	LoadDatabase(file)
	LoadJwtKey(file)
	LoadCache(file)
}
