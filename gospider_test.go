package gospider_test

import (
	"github.com/qianlidongfeng/gospider"
	"testing"
)

func TestGsData_Clone(t *testing.T) {
	var d gospider.Data
	d = gospider.NewGsData().Clone()
	d.Set("hehe",true)
	a:=gospider.NewGsData()
	a.Set("key",123)
	d.Set("xixi",a)
	v,ok:=d.Get("xixi")
	if ok{
		vv,ok:=v.Data().Get("key")
		if ok{
			vvv:=vv.GetValue()
			vvvv:=vvv.(int)
			_=vvvv
		}
	}
}


