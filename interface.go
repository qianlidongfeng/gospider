package gospider

type Data interface{
	Set(string,interface{})
	Get(string) (Value,bool)
	Clone() Data
}

type Value interface{
	Int() int
	Float() float64
	String() string
	Bool() bool
	IntArray() []int
	FloatArray() []float64
	StringArray() []string
	GetValue() interface{}
	Data() Data
}
