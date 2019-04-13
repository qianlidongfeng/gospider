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
	action:=NewAction("homepage","www.baidu.com")
	action.Meta=meta
	action.respy=11
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

