package gospider_test

import (
	"fmt"
	"github.com/qianlidongfeng/gospider"
	"testing"
)

func TestGsData_Clone(t *testing.T) {
	var d gospider.Data
	d = gospider.NewGsData()
	d.Set("bool",true)
	d.Set("slice",[]int{1,2,3})
	a:=gospider.NewGsData()
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
	v.Data().Set("key",456)
	v1,_:=c.Get("xixi")
	i1,_:=v1.Data().Get("key")
	v2,_:=d.Get("xixi")
	i2,_:=v2.Data().Get("key")
	fmt.Println(i1.Int())
	fmt.Println(i2.Int())
}

func TestMeta_Clone(t *testing.T) {
	d := gospider.NewMeta()
	d["bool"]=true
	d["slice"]=[]int{1,2,3}
	a:=gospider.NewMeta()
	a["key"]=123
	d["xixi"]=a
	c:=d.Clone()
	c["bool"]=false
	sl,_:=c["slice"]
	sl.([]int)[0]=5
	fmt.Println(d)
	fmt.Println(c)
	v,ok:=c["xixi"]
	if !ok{
		t.Error("Set key failed")
	}
	v.(gospider.Meta)["key"]=456
	v1,_:=c["xixi"]
	i1,_:=v1.(gospider.Meta)["key"]
	v2,_:=d["xixi"]
	i2,_:=v2.(gospider.Meta)["key"]
	fmt.Println(i1.(int))
	fmt.Println(i2.(int))
}


