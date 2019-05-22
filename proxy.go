package gospider

import (
	"encoding/json"
	"github.com/qianlidongfeng/httpclient"
	"sync"
	"time"
)


type ProxyPool struct{
	Pool chan string
	ProxyServer string
	Size int
	mu sync.Mutex
}

func NewProxyPool(poolSize int,proxyServer string) ProxyPool{
	p:=ProxyPool{}
	p.Size=poolSize
	p.ProxyServer=proxyServer
	p.Pool = make(chan string,poolSize)
	p.mu = sync.Mutex{}
	return p
}

func (this *ProxyPool) Get() string{
	for{
		select{
		case proxy:= <-this.Pool:
			return proxy
		default:
			this.updateProxy()
		}
	}
}

func (this *ProxyPool) updateProxy(){
	this.mu.Lock()
	defer this.mu.Unlock()
	if len(this.Pool) > 0{
		return
	}
	for{
		client:=httpclient.NewHttpClient()
		Resp,err:=client.Get(this.ProxyServer)
		if err != nil{
			log.Warn(err)
			time.Sleep(time.Second*5)
			continue
		}
		var proxies []string
		err = json.Unmarshal([]byte(Resp.Html),&proxies)
		if err != nil{
			log.Warn(err)
			time.Sleep(time.Second*5)
			continue
		}
		if len(proxies)==0{
			log.Msg("updateProxy","update proxy return empty")
			time.Sleep(time.Second*5)
			continue
		}
		for i:=0;i<len(proxies);i++{
			this.Pool<-string(proxies[i])
		}
		break
	}
}