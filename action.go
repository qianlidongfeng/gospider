package gospider

type Parser func (html string,meta Meta,Branch Meta) ([]Action,error)

type Action struct{
	Parser string
	Url string
	Meta Meta
	Branch Meta
	Method string
	PostData string
	failCount int
	respy int
}

func NewAction(parser string,url string) Action{
	return Action{
		Parser:parser,
		Url:url,
		Meta:NewMeta(),
		Branch:NewMeta(),
		Method:"GET",
		PostData:"",
	}
}

func (this *Action) Clone() Action{
	return Action{
		Parser:this.Parser,
		Url:this.Url,
		Meta:this.Meta.Clone(),
		Branch:this.Branch.Clone(),
		Method:this.Method,
		PostData:this.PostData,
	}
}