package gospider

import (
	"errors"
	"github.com/qianlidongfeng/toolbox"
	"gopkg.in/ini.v1"
	"time"
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
	LogerConfig toolbox.MySqlConfig
	EnableCookie bool
	SmartCookie bool
	EnableProxy bool
	EnableDB bool
	ChangeProxy bool
	ChangeAgent bool
	ProxyPoolSize int
	ProxyServer string
	ProxyType string
	TimeOut time.Duration
	ReqClose bool
	ResetHttpclient bool
	DBC DBConfig
	ARC ActionRecordConfig
	Debug bool
	ActionRecord bool
	Delay time.Duration
}

func NewConfig() Config{
	config:=Config{
		LogerConfig:toolbox.MySqlConfig{},
		DBC:DBConfig{},
		ARC:ActionRecordConfig{},
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
		return errors.New("bad ini,spider->thread")
	}
	this.MaxAction,err = s.Key("max_action").Int()
	if err != nil{
		return errors.New("bad ini,spider->max_action")
	}
	this.Debug,err = s.Key("debug").Bool()
	if err != nil{
		return errors.New("bad ini,spider->debug")
	}
	this.EnableCookie,err = s.Key("enable_cookie").Bool()
	if err != nil{
		return errors.New("bad ini,spider->enable_cookie")
	}
	this.SmartCookie,err = s.Key("smart_cookie").Bool()
	if err != nil{
		return errors.New("bad ini,spider->smart_cookie")
	}
	this.EnableProxy,err = s.Key("enable_proxy").Bool()
	if err != nil{
		return errors.New("bad ini,spider->enable_proxy")
	}
	this.ChangeProxy,err = s.Key("change_proxy").Bool()
	if err != nil{
		return errors.New("bad ini,spider->change_proxy")
	}
	this.ProxyPoolSize,err = s.Key("proxypool_size").Int()
	if err != nil{
		return errors.New("bad ini,spider->proxypool_size")
	}
	this.ChangeAgent,err = s.Key("change_agent").Bool()
	if err != nil{
		return errors.New("bad ini,spider->change_agent")
	}
	this.ProxyServer = s.Key("proxyserver").String()
	this.ProxyType = s.Key("proxytype").String()
	timeout,err := s.Key("timeout").Int()
	if err != nil{
		return errors.New("bad ini,spider->timeout")
	}
	this.TimeOut=time.Duration(timeout)*time.Millisecond
	this.ReqClose,err= s.Key("reqclose").Bool()
	if err != nil{
		return errors.New("bad ini,spider->reqclose")
	}
	delay,err := s.Key("delay").Int()
	if err != nil{
		return errors.New("bad ini,spider->delay")
	}
	this.Delay=time.Duration(delay)*time.Millisecond
	this.ResetHttpclient,err = s.Key("reset_httpclient").Bool()
	if err != nil{
		return errors.New("bad ini,spider->reset_httpclient")
	}
	this.EnableDB,err = s.Key("enable_db").Bool()
	if err != nil{
		return errors.New("bad ini,spider->enable_db")
	}
	this.ActionRecord,err = s.Key("action_record").Bool()
	if err != nil{
		return errors.New("bad ini,spider->action_record")
	}
	if this.ActionRecord{
		s= ini.Section("action")
		if s==nil{
			return errors.New("bad ini,cant find action section")
		}
		this.ARC.User = s.Key("user").String()
		this.ARC.PassWord = s.Key("passwd").String()
		this.ARC.DB = s.Key("database").String()
		if this.ARC.DB == ""{
			return errors.New("bad ini,action->database")
		}
		this.ARC.Table = s.Key("table").String()
		if this.ARC.Table == ""{
			return errors.New("bad ini,action->table")
		}
		this.ARC.Type = s.Key("type").String()
		if this.ARC.Type == ""{
			return errors.New("bad ini,action->type")
		}
		this.ARC.Address = s.Key("address").String()
		if this.ARC.Address == ""{
			return errors.New("bad ini,action->address")
		}
		this.ARC.MaxOpenConns,err= s.Key("max_open_conns").Int()
		if err !=nil{
			return errors.New("bad ini,action->max_open_conns")
		}
		this.ARC.MaxIdleConns,err= s.Key("max_idle_conns").Int()
		if err !=nil{
			return errors.New("bad ini,action->max_idle_conns")
		}
		this.ARC.MaxFail,err= s.Key("max_fail").Int()
		if err !=nil{
			return errors.New("bad ini,action->max_respy")
		}
		this.ARC.MaxRespy,err= s.Key("max_respy").Int()
		if err !=nil{
			return errors.New("bad ini,action->max_respy")
		}
		this.ARC.Label= s.Key("label").String()
		if err !=nil || this.ARC.Label==""{
			return errors.New("bad ini,action->label")
		}
	}
	s = ini.Section("loger")
	if s == nil {
		return errors.New("bad ini,cant find loger section")
	}
	this.LogerType = s.Key("type").String()
	if this.LogerType=="netloger"{
		s= ini.Section("logerdb")
		if s == nil{
			return errors.New("bad ini,cant find logerdb section")
		}
		this.LogerConfig.User = s.Key("user").String()
		this.LogerConfig.PassWord = s.Key("passwd").String()
		this.LogerConfig.Address = s.Key("address").String()
		if this.LogerConfig.Address == ""{
			return errors.New("bad ini,logerdb->address")
		}
		this.LogerConfig.DataBase = s.Key("database").String()
		if this.LogerConfig.DataBase == ""{
			return errors.New("bad ini,logerdb->database")
		}
		this.LogerConfig.Table = s.Key("table").String()
		if this.LogerConfig.Table == ""{
			return errors.New("bad ini,ogerdb->table")
		}
		this.LogerConfig.MaxOpenConns,err= s.Key("max_open_conns").Int()
		if err !=nil{
			return errors.New("bad ini,logerdb->max_open_conns")
		}
		this.LogerConfig.MaxIdleConns,err= s.Key("max_idle_conns").Int()
		if err !=nil{
			return errors.New("bad ini,logerdb->max_idle_conns")
		}
	}
	if this.EnableDB{
		s = ini.Section("db")
		if s==nil{
			return errors.New("bad ini,cant find db section")
		}
		this.DBC.User = s.Key("user").String()
		this.DBC.PassWord = s.Key("passwd").String()
		this.DBC.DB = s.Key("database").String()
		if this.DBC.DB == ""{
			return errors.New("bad ini,db->database")
		}
		this.DBC.Type = s.Key("type").String()
		if this.DBC.Type == ""{
			return errors.New("bad ini,db->type")
		}
		this.DBC.Address = s.Key("address").String()
		if this.DBC.Address == ""{
			return errors.New("bad ini,db->address")
		}
		this.DBC.MaxOpenConns,err= s.Key("max_open_conns").Int()
		if err !=nil{
			return errors.New("bad ini,db->max_open_conns")
		}
		this.DBC.MaxIdleConns,err= s.Key("max_idle_conns").Int()
		if err !=nil{
			return errors.New("bad ini,db->max_idle_conns")
		}
	}
	return nil
}

type ActionRecordConfig struct{
	User string
	PassWord string
	Address string
	Type string
	DB string
	Table string
	MaxOpenConns int
	MaxIdleConns int
	MaxRespy int
	MaxFail int
	Label string
}