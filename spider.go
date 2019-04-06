package gospider

import (
	"github.com/qianlidongfeng/httpclient"
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
	pool chan Action
	Clients []httpclient.HttpClient
	opener chan struct{}
	finish bool
	//根据配置文件初始化cd
	//数据库错误日志保存器
	//退出时序列化保存
	//失败时失败actor保存
}

func NewSpider() Spider{
	return Spider{
	}
}

func (this *Spider)Init() error{


	return nil
}

func (this *Spider) Run(){
 	for i:=0;i<len(this.pool);i++{
 		go func(i int){
 			for !this.finish{
				<-this.opener
				var action Action
				action = <-this.pool
				var html string
				var err error
				if strings.ToUpper(action.Method) == "GET"{
					html,err=this.Clients[i].Get(action.Url)
				}else if strings.ToUpper(action.Method) == "POST"{
					html,err=this.Clients[i].Post(action.Url,action.PostData)
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
}

func (this *Spider) AddEntry(action Action){
	this.pool<-action
}

