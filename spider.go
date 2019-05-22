package gospider

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/qianlidongfeng/httpclient"
	"github.com/qianlidongfeng/loger"
	"github.com/qianlidongfeng/loger/netloger"
	syslog "log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)


type savehandle func(meta Meta,db *sql.DB) error

type Spider struct{
	actionPool chan Action
	clients []httpclient.HttpClient
	termsignal bool
	finish bool
	cfg Config
	output *os.File
	onsave savehandle
	db *sql.DB
	actionRecorder Recorder
	parsers map[string]Parser
	mu sync.Mutex
	wg sync.WaitGroup
	threadWg sync.WaitGroup
	gracefulQuitComplete chan struct{}
	mode string
	proxyPool ProxyPool
}

func NewSpider() Spider{
	return Spider{
		parsers:make(map[string]Parser),
		mu:sync.Mutex{},
		wg:sync.WaitGroup{},
		threadWg:sync.WaitGroup{},
		gracefulQuitComplete:make(chan struct{}),
	}
}

func (this *Spider)Init() error{
	this.cfg= NewConfig()
	exe,err:=os.Readlink("/proc/self/exe")
	if err != nil{
		return err
	}
	defaultIni:=exe+".ini"
	configFile:=flag.String("c", defaultIni, "the path of config file")
	mode:=flag.String("m", "run", "the spider mode\nrun:begin\nfix:fix\n")
	flag.Parse()
	this.mode=*mode
	_,err=os.Stat(defaultIni)
	if err != nil{
		return errors.New("spider config file not found")
	}
	err=this.cfg.Init(*configFile)
	if err != nil{
		return err
	}
	//重定向输出流
	if this.cfg.Debug == false{
		this.output, err = os.OpenFile(exe+".log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
		if err != nil{
			return err
		}
		syslog.SetOutput(this.output)
	}
	//初始化日志生成器
	if this.cfg.LogerType=="netloger" && this.cfg.Debug==false{
		lg:=netloger.NewSqloger()
		err=lg.Init(this.cfg.LogerConfig)
		if err != nil{
			return err
		}//
		log=lg
	}else{
		log=loger.NewLocalLoger()
	}
	//初始化数据库
	if this.cfg.EnableDB{
		this.db,err = sql.Open(this.cfg.DBC.Type,this.cfg.DBC.User+":"+
			this.cfg.DBC.PassWord+"@tcp("+this.cfg.DBC.Address+")/"+this.cfg.DBC.DB)
		if err != nil{
			log.Fatal(err)
			return err
		}
		err=this.db.Ping()
		if err != nil{
			log.Fatal(err)
			return err
		}
		this.db.SetMaxOpenConns(this.cfg.DBC.MaxOpenConns)
		this.db.SetMaxIdleConns(this.cfg.DBC.MaxIdleConns)
	}
	//初始化action记录器
	if this.cfg.ActionRecord{
		var label string
		switch runtime.GOOS {
		case "windows":
			label = exe[strings.LastIndex(exe,"\\")+1:]
		default:
			label = exe[strings.LastIndex(exe,"/")+1:]
		}
		this.actionRecorder=&ActionRecorder{
			label:label,
		}
		err=this.actionRecorder.Init(this.cfg.ARC)
		if err != nil{
			return err
		}
	}
	//初始化代理
	this.proxyPool=NewProxyPool(this.cfg.ProxyPoolSize,this.cfg.ProxyServer)
	//初始化线程
	this.actionPool=make(chan Action,this.cfg.MaxAction)
	for i:=0;i<this.cfg.Thread;i++{
		client:=this.NewClient()
		this.clients=append(this.clients,client)
	}
	//启动优雅退出处理线程
	go this.GracefulQuit()
	return nil
}

func (this *Spider) Run(){
 	for i:=0;i<len(this.clients);i++{
 		this.threadWg.Add(1)
 		go func(i int){
 			for !this.termsignal && !this.finish{
				this.Do(i)
			}
 			this.threadWg.Done()
		}(i)
	}
	this.wg.Wait()
 	this.finish=true
 	this.threadWg.Wait()
 	if this.termsignal{
		<-this.gracefulQuitComplete
		os.Exit(1)
	}
	runtime.GC()
}

func (this *Spider) Fix(){
	this.actionPool=make(chan Action,this.cfg.MaxAction)
	actions,err:=this.actionRecorder.GetActions()
	if err != nil{
		log.Fatal(err)
	}
	for _,action:=range actions{
		this.actionPool<-action
	}
	this.Run()
}

func (this *Spider) Start(){
	if this.mode=="fix"{
		fmt.Println("fix")
		this.Fix()
	}else if this.mode == "run"{
		fmt.Println("run")
		this.Run()
	}
}

func (this *Spider) Do(i int){
	var action Action
	action=<-this.actionPool
	defer this.wg.Done()
	if this.cfg.Delay != 0{
		time.Sleep(this.cfg.Delay)
	}
	if this.cfg.Debug{
		log.(*loger.LocalLoger).Debug(action.Method+" "+action.Url)
	}
	var resp httpclient.Resp
	var err error
	if this.cfg.EnableCookie && action.Cookies != nil{
		this.clients[i].SetCookies(action.Url,action.Cookies)
	}
	for k,v :=range action.TempHeader{
		this.clients[i].SetTempHeaderField(k,v)
	}
	if strings.ToUpper(action.Method) == "GET"{
		resp,err=this.clients[i].Get(action.Url)
	}else if strings.ToUpper(action.Method) == "POST"{
		resp,err=this.clients[i].Post(action.Url,action.PostData)
	}
	if err != nil{
		action.failCount++
		if action.failCount>this.cfg.ARC.MaxFail && this.cfg.ActionRecord{
			action.Respy++
			cookies,err:=this.clients[i].GetCooikes(action.Url)
			action.SetCookies(cookies)
			err=this.actionRecorder.Put(action)
			if err != nil{
				log.Warn(err)
			}
		}else{
			this.AddAction(action)
		}
		if this.cfg.EnableProxy&&this.cfg.ChangeProxy{
			if this.cfg.ProxyType=="http"{
				this.clients[i].SetHttpProxy("http://"+this.proxyPool.Get())
			}else if this.cfg.ProxyType=="sock5"{
				this.clients[i].SetSock5Proxy(this.proxyPool.Get())
			}
		}
		log.Warn(err)
		return
	}
	log.Msg("success",action.Url)
	err = this.parsers[action.Parser](this,resp,action.Meta)
	if err != nil {
		log.Warn(err)
		action.failCount++
		if action.failCount>this.cfg.ARC.MaxFail && this.cfg.ActionRecord{
			action.Respy++
			cookies,err:=this.clients[i].GetCooikes(action.Url)
			action.SetCookies(cookies)
			err=this.actionRecorder.Put(action)
			if err != nil{
				log.Warn(err)
			}
		}else{
			this.AddAction(action)
		}
		if this.cfg.EnableProxy&&this.cfg.ChangeProxy{
			if this.cfg.ProxyType=="http"{
				this.clients[i].SetHttpProxy("http://"+this.proxyPool.Get())
			}else if this.cfg.ProxyType=="sock5"{
				this.clients[i].SetSock5Proxy(this.proxyPool.Get())
			}
		}
		return
	}
	if this.cfg.EnableProxy&&this.cfg.ChangeProxy{
		if this.cfg.ProxyType=="http"{
			this.clients[i].SetHttpProxy("http://"+this.proxyPool.Get())
		}else if this.cfg.ProxyType=="sock5"{
			this.clients[i].SetSock5Proxy(this.proxyPool.Get())
		}
	}
	if this.cfg.ChangeAgent{
		this.clients[i].SetHeaderField("User-Agent", httpclient.UserAgents.One())
	}
	if this.cfg.Debug{
		log.(*loger.LocalLoger).Debug(action.Url+" done")
	}
}

func (this *Spider) AddAction(action Action){
	this.wg.Add(1)
	this.actionPool<-action
}

func (this *Spider) Release(){
	if _,ok:=log.(*loger.LocalLoger);ok{
		this.output.Close()
	}
	log.Close()
	if this.cfg.ActionRecord{
		this.actionRecorder.Close()
	}
	if this.cfg.EnableDB{
		this.db.Close()
	}
}

func (this *Spider) SetSaveHander(f savehandle){
	this.onsave=f
}

func (this *Spider) RegisterParser(name string,parser Parser){
	this.mu.Lock()
	this.parsers[name]=parser
	this.mu.Unlock()
}

func (this *Spider) GracefulQuit(){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt,os.Kill,syscall.SIGTERM)
	<-c
	this.termsignal=true
	this.threadWg.Wait()
	record:
	for{
		select{
		case action:=<-this.actionPool:
			if this.cfg.ActionRecord{
				err:=this.actionRecorder.Put(action)
				if err != nil{
					log.Warn(err)
				}
			}
			this.wg.Done()
		default:
			break record
		}
	}
	this.gracefulQuitComplete<-struct{}{}
}

func (this *Spider) SetActionLabel(label string){
	this.actionRecorder.SetActionLabel(label)
}

func (this *Spider) Save(meta Meta){
	err := this.onsave(meta,this.db)
	if err != nil{
		log.Warn(err)
	}
}

func (this *Spider) NewClient() httpclient.HttpClient{
	client:=httpclient.NewHttpClient()
	if this.cfg.TimeOut != 0{
		client.SetTimeOut(this.cfg.TimeOut)
	}
	if this.cfg.EnableCookie{
		client.EnableCookie()
	}
	if this.cfg.EnableProxy{
		if this.cfg.ProxyType=="http"{
			client.SetHttpProxy("http://"+this.proxyPool.Get())
		}else if this.cfg.ProxyType=="sock5"{
			client.SetSock5Proxy(this.proxyPool.Get())
		}
	}
	return client
}

func (this *Spider) Log(label string,v ...interface{}){
	log.Msg(label,v...)
}