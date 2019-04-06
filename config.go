package gospider

import (
	"github.com/headzoo/surf/errors"
	"github.com/qianlidongfeng/loger/netloger"
	"gopkg.in/ini.v1"
)

type DBConfig struct{
	User string
	PassWord string
	Address string
	DB string
	Type string
	MaxOpenConns int
	MaxIdleConns int
}

type Config struct{
	Thread int
	MaxAction int
	LogerType string
	LogerConfig netloger.SqlConfig
	DBC DBConfig
	Debug bool
}

func NewConfig() Config{
	config:=Config{
		LogerConfig:netloger.SqlConfig{},
		DBC:DBConfig{},
	}
	return config
}

func (this *Config) Init(configFile string) error{
	ini, err := ini.Load(configFile)
	if err != nil{
		return err
	}
	s := ini.Section("spider")
	if s == nil {
		return errors.New("bad ini,cant find spider section")
	}
	this.Thread,err=s.Key("thread").Int()
	if err != nil{
		return errors.New(err.Error(),"thread")
	}
	this.MaxAction,err = s.Key("max_action").Int()
	if err != nil{
		return errors.New(err.Error(),"max_action")
	}
	this.Debug,err = s.Key("debug").Bool()
	if err != nil{
		this.Debug = true
	}
	s = ini.Section("loger")
	if s == nil {
		return errors.New("bad ini,cant find loger section")
	}
	this.LogerType = s.Key("type").String()
	if this.LogerType=="netloger"{
		s= ini.Section("logerdb")
		if s == nil{
			return errors.New("bad ini,netloger cant find logerdb section")
		}
		this.LogerConfig.User = s.Key("user").String()
		this.LogerConfig.PassWord = s.Key("passwd").String()
		this.LogerConfig.Address = s.Key("address").String()
		if this.LogerConfig.Address == ""{
			return errors.New("bad ini,netloger loger->address is empty")
		}
		this.LogerConfig.Type = s.Key("dbtype").String()
		if this.LogerConfig.Type == ""{
			return errors.New("bad ini,netloger loger->dbtype is empty")
		}
		this.LogerConfig.DB = s.Key("database").String()
		if this.LogerConfig.DB == ""{
			return errors.New("bad ini,netloger loger->database is empty")
		}
		this.LogerConfig.Table = s.Key("table").String()
		if this.LogerConfig.Table == ""{
			return errors.New("bad ini,netloger loger->table is empty")
		}
		this.LogerConfig.MaxOpenConns,err= s.Key("max_open_conns").Int()
		if err !=nil{
			return errors.New("bad ini,netloger loger->max_open_conns error")
		}
		this.LogerConfig.MaxIdleConns,err= s.Key("max_idle_conns").Int()
		if err !=nil{
			return errors.New("bad ini,netloger loger->max_idle_conns error")
		}
	}
	s = ini.Section("db")
	if s==nil{
		return errors.New("bad ini,cant find db section")
	}
	this.DBC.User = s.Key("user").String()
	this.DBC.PassWord = s.Key("passwd").String()
	this.DBC.DB = s.Key("database").String()
	if this.DBC.DB == ""{
		return errors.New("bad ini,db->database is empty")
	}
	this.DBC.Type = s.Key("type").String()
	if this.DBC.Type == ""{
		return errors.New("bad ini,db->type is empty")
	}
	this.DBC.Address = s.Key("address").String()
	if this.DBC.Address == ""{
		return errors.New("bad ini,db->address is empty")
	}
	this.DBC.MaxOpenConns,err= s.Key("max_open_conns").Int()
	if err !=nil{
		return errors.New("bad ini,netloger loger->max_open_conns error")
	}
	this.DBC.MaxIdleConns,err= s.Key("max_idle_conns").Int()
	if err !=nil{
		return errors.New("bad ini,netloger loger->max_idle_conns error")
	}
	return nil
}