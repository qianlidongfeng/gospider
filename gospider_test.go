package gospider_test

import (
	"fmt"
	"github.com/qianlidongfeng/gospider"
	"testing"
)

func TestMeta_Clone(t *testing.T) {
	d := gospider.NewMeta()
	d.Set("bool",true)
	d.Set("slice",[]int{1,2,3})
	a:=gospider.NewMeta()
	a.Set("key",123)
	d.Set("xixi",a)
	c:=d.Clone()
	c.Set("bool",false)
	sl,_:=c.Get("slice")
	sl.IntArray()[0]=5
	fmt.Println(d)
	fmt.Println(c)
	v,ok:=c.Get("xixi")
	if !ok{
		t.Error("Set key failed")
	}
	v.Meta().Set("key",456)
	v1,_:=c.Get("xixi")
	i1,_:=v1.Meta().Get("key")
	v2,_:=d.Get("xixi")
	i2,_:=v2.Meta().Get("key")
	fmt.Println(i1.Int())
	fmt.Println(i2.Int())
	_,ok=d.Get("xixi")
	fmt.Println(ok)
	d.Delete("xixi")
	fmt.Println(d)
	_,ok=d.Get("xixi")
	fmt.Println(ok)
	d.Clear()
	fmt.Println(d)
	fmt.Println(c)
	_,ok=d.Get("xixi")
	fmt.Println(ok)
}

func TestSpider_Init(t *testing.T) {
	sp:=gospider.NewSpider()
	err:=sp.Init()
	if err != nil{
		t.Error(err)
	}
}

func TestAll(t *testing.T) {
	sp:=gospider.NewSpider()
	err:=sp.Init()
	if err != nil{
		t.Error(err)
	}
}


