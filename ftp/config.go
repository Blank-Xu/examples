package ftp

type Config struct {
	ServerName  string `json:"server_name" xml:"server_name" yaml:"server_name" toml:"server_name"`
	Dir         string `json:"dir" xml:"dir" yaml:"dir" toml:"dir"`
	Host        string `json:"host" xml:"host" yaml:"host" toml:"host"`
	Port        int    `json:"port" xml:"port" yaml:"port" toml:"port"`
	PasvMinPort uint   `json:"pasv_min_port" xml:"pasv_min_port" yaml:"pasv_min_port" toml:"pasv_min_port"`
	PasvMaxPort uint   `json:"pasv_max_port" xml:"pasv_max_port" yaml:"pasv_max_port" toml:"pasv_max_port"`

	TlsKey  string
	AutoTls bool

	Users *struct {
		Username string
		Password string
	} `json:"users" xml:"users" yaml:"users" toml:"users"`

	addr string
}

func (p *Config) init() {
	if len(p.ServerName) == 0 {
		p.ServerName = "FTP Server"
	}
	if p.Port == 0 {
		p.Port = 21
	}

	p.addr = GetTcpAddr(p.Host, p.Port)
}
