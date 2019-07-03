package gospider

import (
	"github.com/qianlidongfeng/httpclient"
	"net/http"
)

type Parser func (spider *Spider,resp httpclient.Resp,meta Meta) error

type Action struct{
	Parser string
	Url string
	Meta Meta
	Method string
	PostData []byte
	TempHeader map[string]string
	Cookies []*http.Cookie
	failCount int
	Respy int
}

func NewAction() Action{
	return Action{
		Meta:NewMeta(),
		Method:"GET",
		PostData:nil,
		TempHeader:make(map[string]string),
	}
}

func (this *Action) Clone() Action{
	a:= Action{
		Parser:this.Parser,
		Url:this.Url,
		Meta:this.Meta.Clone(),
		Method:this.Method,
		PostData:this.PostData,
	}
	a.TempHeader=make(map[string]string)
	for k,v:=range this.TempHeader{
		a.TempHeader[k]=v
	}
	a.SetCookies(this.Cookies)
	return a
}

func (this *Action) UnsafeClone() Action{
	a:= Action{
		Parser:this.Parser,
		Url:this.Url,
		Meta:this.Meta.UnsafeClone(),
		Method:this.Method,
		PostData:this.PostData,
	}
	a.TempHeader=make(map[string]string)
	for k,v:=range this.TempHeader{
		a.TempHeader[k]=v
	}
	a.SetCookies(this.Cookies)
	return a
}

func (this *Action) SetParser(parser string){
	this.Parser=parser
}

func (this *Action) SetUrl(url string){
	this.Url=url
}

func (this *Action) SetMeta(meta Meta){
	this.Meta=meta.Clone()
}

func (this *Action) SetMethod(method string){
	this.Method=method
}

func (this *Action) SetPostData(postdata []byte){
	this.PostData=postdata
}

func (this *Action) SetTempHeaderField(key string,value string){
	this.TempHeader[key]=value
}

func (this *Action) SetCookiesByString(domain string,path string,cookie string) error{
	var err error
	this.Cookies,err=httpclient.MakeCookies(domain,path,cookie)
	if err != nil{
		return err
	}
	return nil
}

func (this *Action) SetCookies(cookies []*http.Cookie){
	this.Cookies=cookies
}