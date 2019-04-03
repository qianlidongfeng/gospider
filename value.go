package gospider

type gsValue struct{
	value interface{}
}

func (this *gsValue) Int() int{
	return this.value.(int)
}

func (this *gsValue) Float() float64{
	return this.value.(float64)
}

func (this *gsValue) String() string{
	return this.value.(string)
}

func (this *gsValue) Bool() bool{
	return this.value.(bool)
}

func (this *gsValue) IntArray() []int{
	return this.value.([]int)
}

func (this *gsValue) FloatArray() []float64{
	return this.value.([]float64)
}

func (this *gsValue) StringArray() []string{
	return this.value.([]string)
}

func (this *gsValue) GetValue() interface{}{
	return this.value
}

func (this *gsValue) Data() Data{
	return this.value.(Data)
}