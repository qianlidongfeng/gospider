package gospider

import (
	"encoding/gob"
	"fmt"
	"os"
	"testing"
)



func TestAction_Clone(t *testing.T) {
	meta:=NewMeta()
	meta.Set("int",1)
	meta.Set("str","fuck")
	meta.Set("slice",[]string{"yang","huan","anal"})
	data:=NewMeta()
	data.Set("one","wocao")
	data.Set("two","nima")
	data.Set("array",[]string{"haha","hehe","heihei"})
	meta.Set("data",data)
	var ms []Meta
	ms=append(ms,meta.Clone())
	ms=append(ms,meta.Clone())
	ms=append(ms,meta.Clone())
	ms=append(ms,meta.Clone())
	file, err := os.Create("gob")
	if err != nil {
		fmt.Println(err)
	}
	gob.Register(Meta{})
	enc := gob.NewEncoder(file)
	err2 := enc.Encode(meta)
	file.Close()
	fmt.Println(err2)
	sp:=NewSpider()
	err=sp.Init()
	if err != nil{
		t.Error(err)
	}
	action:=NewAction()
	action.Meta=meta
	action.Respy=11
	sp.actionRecorder.Put(action)
	var ms2 []Meta
	var meta2 Meta
	actions,err:=sp.actionRecorder.GetActions()
	_=actions
	_=ms2
	_=meta2
	d:=gob.NewDecoder(file)
	err=d.Decode(&meta2)
	fmt.Println(err)
	file.Close()
}

func TestNewProxyPool(t *testing.T) {
	pt:=NewProxyPool(100,"http://127.0.0.1:8080/httpproxies?count=10")
	for{
		proxy:=pt.Get()
		fmt.Println(proxy)
	}
}

func TestRun(t *testing.T) {
	sp:=NewSpider()
	err:=sp.Init()
	if err != nil{
		t.Error(err)
	}
	action:=NewAction()
	action.SetCookiesByString(".baidu.com","/",`BAIDUID=902ACFD6E02BB6B47AD388A5A58921C3:FG=1; PSTM=1555830874; BD_UPN=123353; BIDUPSID=BF28140F3AACB62E8B2DAF85E3F51CAB; BDUSS=5FdHpPUVVscXp6bGVobjJjUXFocUtqVDc4UE1VV2xzTGFBTEEyQi1UeUZyT05jRVFBQUFBJCQAAAAAAAAAAAEAAAA4vbUdQW5nbGVfU2VhbgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIUfvFyFH7xcTG; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; BD_HOME=1; H_PS_PSSID=28959_1457_28938_21103_28768_28724_28964_28832_28585`)
	action.SetUrl("https://www.baidu.com")
	action.SetParser("home")
	action.SetTempHeaderField("referer","https://www.tieba.com")
	d := NewMeta()
	d.Set("bool",true)
	d.Set("slice",[]int{1,2,3})
	a:=NewMeta()
	a.Set("key",123)
	d.Set("xixi",a)
	action.SetMeta(d)
	action.Respy=8
	sp.actionRecorder.Put(action)
	sp.Fix()
}