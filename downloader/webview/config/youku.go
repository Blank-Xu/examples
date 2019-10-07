package config

type Youku struct {
	Ccode    string `yaml:"ccode"`
	Ckey     string `yaml:"ckey"`
	Password string `yaml:"password"`
}

func (p *Youku) Init() {
	if len(p.Ccode) == 0 {

	}
	if len(p.Ckey) == 0 {

	}
	if len(p.Password) == 0 {

	}
}
