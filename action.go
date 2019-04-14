package gospider

type Parser func (spider *Spider,html string,meta Meta) error

type Action struct{
	Parser string
	Url string
	Meta Meta
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
		Method:"GET",
		PostData:"",
	}
}

func (this *Action) Clone() Action{
	return Action{
		Parser:this.Parser,
		Url:this.Url,
		Meta:this.Meta.Clone(),
		Method:this.Method,
		PostData:this.PostData,
	}
}

func (this *Action) SetParser(parser string){
	this.Parser=parser
}

func (this *Action) SetUrl(url string){
	this.Url=url
}

func (this *Action) SetMeta(meta Meta){
	this.Meta=meta.Clone()
}

func (this *Action) SetMethod(method string){
	this.Method=method
}

func (this *Action) SetPostData(postdata string){
	this.PostData=postdata
}