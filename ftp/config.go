package ftp

type Config struct {
	ServerName  string `json:"server_name" xml:"server_name" yaml:"server_name" toml:"server_name"`
	HostName    string `json:"host_name" xml:"host_name" yaml:"host_name" toml:"host_name"`
	Port        int    `json:"port" xml:"port" yaml:"port" toml:"port"`
	PasvMinPort uint   `json:"pasv_min_port" xml:"pasv_min_port" yaml:"pasv_min_port" toml:"pasv_min_port"`
	PasvMaxPort uint   `json:"pasv_max_port" xml:"pasv_max_port" yaml:"pasv_max_port" toml:"pasv_max_port"`

	Ip string `json:"ip" xml:"ip" yaml:"ip" toml:"ip"`
	// Port uint8  `json:"port" xml:"port" yaml:"port" toml:"port"`
	Mode string

	TlsKey  string
	AutoTls bool

	Users *struct {
		Username string
		Password string
	} `json:"users" xml:"users" yaml:"users" toml:"users"`
}

func (p *Config) init() {
	if len(p.ServerName) == 0 {
		p.ServerName = "FTP Server"
	}
	if len(p.HostName) == 0 {
		p.HostName = "127.0.0.1"
	}
	if p.Port == 0 {
		p.Port = 21
	}
}
