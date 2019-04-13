package gospider

type Recorder interface{
	Init(ActionRecordConfig) error
	Put(Action) error
	SetActionLabel(label string)
	GetActions() ([]Action,error)
	Close()
}