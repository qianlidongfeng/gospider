package gospider_test

import (
	"bytes"
	"encoding/gob"
	"github.com/qianlidongfeng/gospider"
	"testing"
)

func TestMeta_Clone(t *testing.T) {
	action:=gospider.NewAction()
	action.SetCookiesByString(".baidu.com","/",`BAIDUID=902ACFD6E02BB6B47AD388A5A58921C3:FG=1; PSTM=1555830874; BD_UPN=123353; BIDUPSID=BF28140F3AACB62E8B2DAF85E3F51CAB; BDUSS=5FdHpPUVVscXp6bGVobjJjUXFocUtqVDc4UE1VV2xzTGFBTEEyQi1UeUZyT05jRVFBQUFBJCQAAAAAAAAAAAEAAAA4vbUdQW5nbGVfU2VhbgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIUfvFyFH7xcTG; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; BD_HOME=1; H_PS_PSSID=28959_1457_28938_21103_28768_28724_28964_28832_28585`)
	action.SetUrl("https://www.baidu.com")
	action.SetParser("home")
	action.SetTempHeaderField("referer","https://www.tieba.com")
	d := gospider.NewMeta()
	d.Set("bool",true)
	d.Set("slice",[]int{1,2,3})
	a:=gospider.NewMeta()
	a.Set("key",123)
	d.Set("xixi",a)
	action.SetMeta(d)
	gob.Register(gospider.Meta{})
	var binary bytes.Buffer
	encoder := gob.NewEncoder(&binary)
	err:=encoder.Encode(action)
	if err != nil{
		t.Error(err)
	}
	b:=binary.Bytes()
	var actiontemp gospider.Action
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	err=decoder.Decode(&actiontemp)
	if err!=nil{
		return
	}
	actionclone:=actiontemp.UnsafeClone()
	_=actionclone
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


