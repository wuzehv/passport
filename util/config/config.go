package config

import (
	"fmt"
	"github.com/go-ini/ini"
	log2 "log"
	"os"
	"time"
)

const (
	Dev = "dev"
	Tra = "tra"
	Pro = "pro"
)

type app struct {
	Env     string
	Port    string
	RootDir string
	Domain  string
	RunMode string
}

type db struct {
	Host               string
	User               string
	Passwd             string
	DbName             string
	Charset            string
	MaxIdleConn        int
	MaxActiveConn      int
	MaxConnIdleTimeout time.Duration
}

type redis struct {
	Host               string
	Passwd             string
	DbNum              int
	MaxIdleConn        int
	MaxActiveConn      int
	MaxConnIdleTimeout time.Duration
}

type log struct {
	Dir      string
	Filename string
}

type svc struct {
	Adapter    string
	ExpireTime time.Duration
}

var (
	config *ini.File
	App    = &app{}
	Db     = &db{}
	Redis  = &redis{}
	Log    = &log{}
	Svc    = &svc{}
)

// ini文件加载优先级
var priorityIni = [...]string{
	"pro.ini",
	"tra.ini",
	"app.ini",
}

func init() {
	var err error
	config, err = ini.Load(getConfigFile())
	if err != nil {
		log2.Fatalf("config load error: %v\n", err)
	}

	mapTo(ini.DefaultFormatLeft, App)
	mapTo("db", Db)
	mapTo("redis", Redis)
	mapTo("log", Log)
	mapTo("svc", Svc)
}

func getConfigFile() string {
	path, err := os.Getwd()
	if err != nil {
		log2.Fatalf("config getpwd error: %v\n", err)
	}

	for _, v := range priorityIni {
		configFile := fmt.Sprintf("%s%sconf%s%s", path, string(os.PathSeparator), string(os.PathSeparator), v)
		if _, err = os.Stat(configFile); !os.IsNotExist(err) {
			return configFile
		}
	}

	log2.Fatalf("config file not exists\n")
	return ""
}

func mapTo(section string, v interface{}) {
	err := config.Section(section).MapTo(v)
	if err != nil {
		log2.Fatalf("config mapto error: %v\n", err)
	}
}

func IsDev() bool {
	return App.Env == Dev
}

func IsTra() bool {
	return App.Env == Tra
}

func IsPro() bool {
	return App.Env == Pro
}
