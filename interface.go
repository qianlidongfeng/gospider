package gospider

type Data interface{
	Set(string,interface{})
	Get(string) (Value,bool)
	Clone() Data
	Length() int
	Delete(key string)
	Clear()
	AddReference() int
	SubReference() int
	GetReference() int
	Map() map[string]interface{}
}

type Value2 interface{
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

type Recorder interface{
	Init(ActionRecordConfig) error
	Put(Action) error
	SetActionLabel(label string)
	GetActions() ([]Action,error)
	Close()
}