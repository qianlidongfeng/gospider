package gospider

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/headzoo/surf/errors"
	"github.com/qianlidongfeng/httpclient"
	"github.com/qianlidongfeng/loger"
	"github.com/qianlidongfeng/loger/netloger"
	"log"
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
	cfg Config
	loger loger.Loger
	output *os.File
	onsave savehandle
	db *sql.DB
	actionRecorder Recorder
	parsers map[string]Parser
	mu sync.Mutex
	wg sync.WaitGroup
	gracefulQuitComplete chan struct{}
	mode string
	//根据配置文件初始化
	//数据库错误日志保存器
	//接口 结构 指针
	//爬虫逻辑
	//Meta是否能反序列化
	//退出时序列化保存
	//失败时失败actor保存
}

func NewSpider() Spider{
	return Spider{
		parsers:make(map[string]Parser),
		mu:sync.Mutex{},
		wg:sync.WaitGroup{},
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
		log.SetOutput(this.output)
	}
	//初始化日志生成器
	if this.cfg.LogerType=="netloger" && this.cfg.Debug==false{
		lg:=netloger.NewSqloger()
		err=lg.Init(this.cfg.LogerConfig)
		if err != nil{
			return err
		}//
		this.loger=lg
	}else{
		this.loger=loger.NewLocalLoger()
	}
	//初始化数据库
	this.db,err = sql.Open(this.cfg.DBC.Type,this.cfg.DBC.User+":"+
		this.cfg.DBC.PassWord+"@tcp("+this.cfg.DBC.Address+")/"+this.cfg.DBC.DB)
	if err != nil{
		this.loger.Fatal(err)
		return err
	}
	err=this.db.Ping()
	if err != nil{
		this.loger.Fatal(err)
		return err
	}
	this.db.SetMaxOpenConns(this.cfg.DBC.MaxOpenConns)
	this.db.SetMaxIdleConns(this.cfg.DBC.MaxIdleConns)
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
	//初始化线程
	this.actionPool=make(chan Action,this.cfg.MaxAction)
	for i:=0;i<this.cfg.Thread;i++{
		client,err:=this.NewClient()
		if err != nil{
			this.loger.Fatal(err)
			return err
		}
		this.clients=append(this.clients,client)
	}
	//启动优雅退出处理线程
	go this.GracefulQuit()
	return nil
}

func (this *Spider) Run(){
 	for i:=0;i<len(this.clients);i++{
 		this.wg.Add(1)
 		go func(i int){
 			var ctn bool = true
 			for !this.termsignal{
				var action Action
				select{
				case action=<-this.actionPool:
					if !ctn{
						ctn=true
						this.wg.Add(1)
					}
					break
				case <-time.After(time.Second):
					if ctn{
						ctn=false
						this.wg.Done()
					}
					continue
				}
				if this.cfg.Delay != 0{
					time.Sleep(this.cfg.Delay)
				}
				if this.cfg.Debug{
					this.loger.(*loger.LocalLoger).Debug(action.Method+" "+action.Url)
				}
				var html string
				var err error
				if strings.ToUpper(action.Method) == "GET"{
					html,err=this.clients[i].Get(action.Url)
				}else if strings.ToUpper(action.Method) == "POST"{
					html,err=this.clients[i].Post(action.Url,action.PostData)
				}
				if err != nil{
					action.failCount++
					if action.failCount>this.cfg.ARC.MaxFail && this.cfg.ActionRecord{
						action.respy++
						err=this.actionRecorder.Put(action)
						if err != nil{
							this.loger.Warn(err)
						}
					}else{
						this.AddAction(action)
					}
					if this.cfg.Debug{
						this.loger.(*loger.LocalLoger).Debug(err)
					}
					if this.cfg.ResetHttpclient{
						this.clients[i],err=this.NewClient()
						if err != nil{
							this.loger.Fatal(err)
						}
					}
					continue
				}
				err = this.parsers[action.Parser](this,html,action.Meta)
				if err != nil {
					this.loger.Warn(err)
					action.failCount++
					if action.failCount>this.cfg.ARC.MaxFail && this.cfg.ActionRecord{
						action.respy++
						err=this.actionRecorder.Put(action)
						if err != nil{
							this.loger.Warn(err)
						}
					}else{
						this.AddAction(action)
					}
					continue
				}
				if this.cfg.Debug{
					this.loger.(*loger.LocalLoger).Debug(action.Url+" done")
				}
			}
 			if ctn{
 				ctn=false
 				this.wg.Done()
			}
		}(i)
	}
	this.wg.Wait()
 	if this.termsignal{
		<-this.gracefulQuitComplete
		this.Release()
	}
	runtime.GC()
}

func (this *Spider) Fix(){
	this.actionPool=make(chan Action,this.cfg.MaxAction)
	actions,err:=this.actionRecorder.GetActions()
	if err != nil{
		this.loger.Fatal()
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

func (this *Spider) AddAction(action Action){
	this.actionPool<-action
}

func (this *Spider) Release(){
	if _,ok:=this.loger.(*loger.LocalLoger);ok{
		this.output.Close()
	}
	this.loger.Close()
	this.actionRecorder.Close()
	this.db.Close()
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
	this.wg.Wait()
	close(this.actionPool)
	for action:=range(this.actionPool){
		err:=this.actionRecorder.Put(action)
		if err != nil{
			this.loger.Warn(err)
		}
	}
	this.gracefulQuitComplete<-struct{}{}
}

func (this *Spider) SetActionLabel(label string){
	this.actionRecorder.SetActionLabel(label)
}

func (this *Spider) Save(meta Meta){
	this.onsave(meta,this.db)
}

func (this *Spider) NewClient() (client httpclient.HttpClient,err error){
	client,err=httpclient.NewHttpClient()
	if err != nil{
		return
	}
	if this.cfg.EnableCookie{
		client.EnableCookie()
	}
	if this.cfg.TimeOut != 0{
		client.SetTimeOut(this.cfg.TimeOut)
	}
	return
}