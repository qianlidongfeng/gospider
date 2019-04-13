package gospider

type Value struct{
	value interface{}
}

func (this *Value) Int() int{
	return this.value.(int)
}

func (this *Value) Float() float64{
	return this.value.(float64)
}

func (this *Value) String() string{
	return this.value.(string)
}

func (this *Value) Bool() bool{
	return this.value.(bool)
}

func (this *Value) IntArray() []int{
	return this.value.([]int)
}

func (this *Value) FloatArray() []float64{
	return this.value.([]float64)
}

func (this *Value) StringArray() []string{
	return this.value.([]string)
}

func (this *Value) GetValue() interface{}{
	return this.value
}

func (this *Value) Meta() *Meta{
	m:=this.value.(Meta)
	return &m
}