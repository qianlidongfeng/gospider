package gospider

import (
	"flag"
	"os"
	"path/filepath"
)

type Config struct{

}

func NewConfig() (config Config,err error){
	config=Config{}
	exe,err:=os.Readlink("/proc/self/exe")
	if err != nil{
		return
	}
	defaultCf:=filepath.Dir(exe)+exe+".ini"
	configFile:=flag.String("c", defaultCf, "config file path")
	flag.Parse()
	config.Init(*configFile)
	return
}

func (this *Config) Init(cfg string) error{
	
	return nil
}