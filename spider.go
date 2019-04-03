package gospider

import (
	"os"
	"os/signal"
	"syscall"
)

type Action struct{
	Parser func (string,Data) (Action,error)
	Urls []string
	Data Data
}


type Spider struct{
	entry Action
	pool chan Action
	opener chan struct{}
	finish bool
}

func NewSpider(maxGoruntime int) Spider{
	return Spider{
		pool:make(chan Action,maxGoruntime),
	}
}


func (this *Spider) AppendEntry(action Action){
	this.entry.Parser = action.Parser
	for _,v := range action.Urls{
		this.entry.Urls=append(this.entry.Urls,v)
	}
}

func (this *Spider) Run(){
 	for i:=0;i<len(this.pool);i++{
 		go func(){
 			for !this.finish{
				<-this.opener
				var action Action
				action = <-this.pool
				for url := range action.Urls {
					_ = url
					html := "abc"
					result, err := action.Parser(html, action.Data.Clone())
					if err != nil {
						continue
					}
					if len(result.Urls) != 0 {
						this.pool <- result
					}else{

					}
				}
			}
		}()
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt,os.Kill,syscall.SIGTERM)
	<-c
}

