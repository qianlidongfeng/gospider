package gospider

import (
	"reflect"
	"sync"
)


type Meta struct{
	M map[string] interface{}
	mu *sync.Mutex
}

func NewMeta() Meta{
	return Meta{
		M:make(map[string]interface{}),
		mu:&sync.Mutex{},
	}
}

func(this *Meta) Set(key string,value interface{}){
	this.mu.Lock()
	defer this.mu.Unlock()
	this.M[key]=value
}

func(this *Meta) Get(key string) (v Value,ok bool){
	this.mu.Lock()
	defer this.mu.Unlock()
	vl,ok:= this.M[key]
	return Value{value:vl},ok
}

func(this *Meta) Clone() Meta{
	this.mu.Lock()
	defer this.mu.Unlock()
	m:=NewMeta()

	for k,v := range this.M{
		rtype := reflect.TypeOf(v)
		switch rtype.Kind() {
		case reflect.Slice:
			if s,ok:=v.([]int);ok {
				s = make([]int, len(v.([]int)))
				copy(s, v.([]int))
				m.Set(k, s)
			}else if s,ok:=v.([]float64);ok{
				s = make([]float64, len(v.([]float64)))
				copy(s, v.([]float64))
				m.Set(k, s)
			}else if s,ok:=v.([]string);ok{
				s = make([]string, len(v.([]string)))
				copy(s, v.([]string))
				m.Set(k, s)
			}else if s,ok:=v.([]interface{});ok{
				s = make([]interface{}, len(v.([]interface{})))
				copy(s, v.([]interface{}))
				m.Set(k, s)
			}else{
				m.Set(k,v)
			}
		case reflect.Struct:
			if s,ok:=v.(Meta);ok{
				s = v.(Meta)
				m.Set(k,s.Clone())
			}else{
				m.Set(k,v)
			}
		default:
			m.Set(k,v)
		}
	}
	return m
}

func (this *Meta) Clear(){
	this.mu.Lock()
	defer this.mu.Unlock()
	this.M=make(map[string]interface{})
}

func (this *Meta) Length() int{
	this.mu.Lock()
	defer this.mu.Unlock()
	return len(this.M)
}

func (this *Meta) Delete(key string){
	this.mu.Lock()
	defer this.mu.Unlock()
	delete(this.M,key)
}

func (this *Meta) Map() map[string]interface{}{
	return this.Clone().M
}