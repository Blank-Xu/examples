package ftp

type Config struct {
	Ip   string `json:"ip" xml:"ip" yaml:"ip" toml:"ip"`
	Port uint8  `json:"port" xml:"port" yaml:"port" toml:"port"`
	Mode string

	TlsKey  string
	AutoTls bool

	Users *struct {
		Username string
		Password string
	} `json:"users" xml:"users" yaml:"users" toml:"users"`
}
