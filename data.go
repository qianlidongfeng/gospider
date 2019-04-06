package gospider

import (
	"reflect"
)

type Meta map[string] interface{}


func NewMeta() Meta{
	return 	make(map[string]interface{})

}

func(this Meta) Clone() Meta{
	d:=NewMeta()
	for k,v := range this{
		rtype := reflect.TypeOf(v)
		switch rtype.Kind() {
		case reflect.Slice:
			if s,ok:=v.([]int);ok {
				s = make([]int, len(v.([]int)))
				copy(s, v.([]int))
				d[k]=s
			}else if s,ok:=v.([]float64);ok{
				s = make([]float64, len(v.([]float64)))
				copy(s, v.([]float64))
				d[k]=s
			}else if s,ok:=v.([]string);ok{
				s = make([]string, len(v.([]string)))
				copy(s, v.([]string))
				d[k]=s
			}else if s,ok:=v.([]interface{});ok{
				s = make([]interface{}, len(v.([]interface{})))
				copy(s, v.([]interface{}))
				d[k]=s
			}else{
				d[k]=v
			}
		case reflect.Map:
			if s,ok:=v.(Meta);ok{
				s = v.(Meta).Clone()
				d[k]=s
			}else{
				d[k]=v
			}
		default:
			d[k]=v
		}
	}
	return d
}

type GsData struct{
	m map[string] interface{}
}

func NewGsData() *GsData{
	return &GsData{
		m:make(map[string]interface{}),
	}
}

func(this *GsData) Set(key string,value interface{}){
	this.m[key]=value
}

func(this *GsData) Get(key string) (v Value,ok bool){
	vl,ok:= this.m[key]
	return &gsValue{value:vl},ok
}

func(this *GsData) Clone() Data{
	d:=NewGsData()
	for k,v := range this.m{
		rtype := reflect.TypeOf(v)
		switch rtype.Kind() {
		case reflect.Slice:
			if s,ok:=v.([]int);ok {
				s = make([]int, len(v.([]int)))
				copy(s, v.([]int))
				d.Set(k, s)
			}else if s,ok:=v.([]float64);ok{
				s = make([]float64, len(v.([]float64)))
				copy(s, v.([]float64))
				d.Set(k, s)
			}else if s,ok:=v.([]string);ok{
				s = make([]string, len(v.([]string)))
				copy(s, v.([]string))
				d.Set(k, s)
			}else if s,ok:=v.([]interface{});ok{
				s = make([]interface{}, len(v.([]interface{})))
				copy(s, v.([]interface{}))
				d.Set(k, s)
			}else{
				d.Set(k,v)
			}
		case reflect.Ptr:
			if s,ok:=v.(Data);ok{
				s = v.(Data).Clone()
				d.Set(k,s)
			}else{
				d.Set(k,v)
			}
		default:
			d.Set(k,v)
		}
	}
	return d
}

func (this *GsData) Clear(){
	this=NewGsData()
}

func (this *GsData) Length() int{
	return len(this.m)
}

func (this *GsData) Delete(key string){
	delete(this.m,key)
}

