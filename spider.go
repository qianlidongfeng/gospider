package gospider

import (
	"flag"
	"github.com/headzoo/surf/errors"
	"github.com/qianlidongfeng/httpclient"
	"github.com/qianlidongfeng/loger"
	"github.com/qianlidongfeng/loger/netloger"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Parser func (sp *Spider,html string,meta Meta,extra Meta) (Action,error)

type Action struct{
	Parser Parser
	Url string
	Meta Meta
	Extra Meta
	Method string
	PostData string
}

func NewAction() Action{
	return Action{}
}

func (this *Action) Clone() Action{
	return Action{
		Parser:this.Parser,
		Url:this.Url,
		Meta:this.Meta.Clone(),
		Extra:this.Extra.Clone(),
		Method:this.Method,
		PostData:this.PostData,
	}
}


type Spider struct{
	actions chan Action
	clients []httpclient.HttpClient
	opener chan struct{}
	finish bool
	cfg Config
	loger loger.Loger
	output *os.File
	//根据配置文件初始化
	//数据库错误日志保存器
	//退出时序列化保存
	//失败时失败actor保存
}

func NewSpider() Spider{
	return Spider{
	}
}

func (this *Spider)Init() error{
	this.cfg= NewConfig()
	exe,err:=os.Readlink("/proc/self/exe")
	if err != nil{
		return err
	}
	defaultIni:=exe+".ini"
	configFile:=flag.String("spc", defaultIni, "the path of config file")
	flag.Parse()
	_,err=os.Stat(defaultIni)
	if err != nil{
		return errors.New("spider config file not found")
	}
	err=this.cfg.Init(*configFile)
	if err != nil{
		return err
	}
	if this.cfg.LogerType=="netloger" && this.cfg.LogerConfig.Type=="mysql" && this.cfg.Debug==false{
		this.loger=netloger.NewSqloger()
		if l,ok:=this.loger.(*netloger.Sqloger);ok{
			err=l.Init(this.cfg.LogerConfig)
			if err != nil{
				return err
			}//if l,ok:=this.loger.(*netloger.Oracleloger);ok
		}else{
			return errors.New("assert netloger.Sqloger failed")
		}
	}else{
		this.loger=loger.NewLocalLoger()
		if this.cfg.Debug == false{
			this.output, err = os.OpenFile(exe+".log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
			if err != nil{
				return err
			}
			this.loger.(*loger.LocalLoger).SetOutPut(this.output)
		}
	}
	this.actions=make(chan Action,this.cfg.MaxAction)
	for i:=0;i<this.cfg.Thread;i++{
		client,err:=httpclient.NewHttpClient()
		if err != nil{
			this.loger.Fatal(err)
			return err
		}
		this.clients=append(this.clients,client)
	}
	return nil
}

func (this *Spider) Run(){
 	for i:=0;i<len(this.actions);i++{
 		go func(i int){
 			for !this.finish{
				<-this.opener
				var action Action
				action = <-this.actions
				var html string
				var err error
				if strings.ToUpper(action.Method) == "GET"{
					html,err=this.clients[i].Get(action.Url)
				}else if strings.ToUpper(action.Method) == "POST"{
					html,err=this.clients[i].Post(action.Url,action.PostData)
				}
				result, err := action.Parser(this,html,action.Meta,action.Extra)
				if err != nil {
					continue
				}
				_=result
			}
		}(i)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt,os.Kill,syscall.SIGTERM)
	<-c
	this.Release()
}

func (this *Spider) AddEntry(action Action){
	this.actions<-action
}

func (this *Spider) Release(){
	if l,ok:=this.loger.(*netloger.Sqloger);ok{
		l.Release()
	}else if _,ok:=this.loger.(*loger.LocalLoger);ok{
		this.output.Close()
	}
}