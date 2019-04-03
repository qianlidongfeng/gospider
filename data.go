package gospider

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
		d.Set(k,v)
	}
	return d
}


