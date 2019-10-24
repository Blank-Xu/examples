package ftp

import (
	"errors"
	"os"
)

type Config struct {
	ServerName             string `json:"server_name" xml:"server_name" yaml:"server_name" toml:"server_name"`
	Dir                    string `json:"dir" xml:"dir" yaml:"dir" toml:"dir"`
	Host                   string `json:"host" xml:"host" yaml:"host" toml:"host"`
	Port                   uint32 `json:"port" xml:"port" yaml:"port" toml:"port"`
	MaxConnections         uint32 `json:"max_connections" xml:"max_connections" yaml:"max_connections" toml:"max_connections"`
	DeadlineSeconds        uint32 `json:"deadline_seconds"`
	ReadDeadlineSeconds    uint32 `json:"read_deadline_seconds"`
	WriteDeadlineSeconds   uint32 `json:"write_deadline_seconds"`
	KeepAlivePeriodSeconds uint32 `json:"keep_alive_period_seconds"`

	PasvMinPort uint32 `json:"pasv_min_port" xml:"pasv_min_port" yaml:"pasv_min_port" toml:"pasv_min_port"`
	PasvMaxPort uint32 `json:"pasv_max_port" xml:"pasv_max_port" yaml:"pasv_max_port" toml:"pasv_max_port"`

	TlsKey  string
	AutoTls bool

	Accounts []*Account `json:"accounts" xml:"accounts" yaml:"accounts" toml:"accounts"`

	addr       string
	accountMap map[string]*Account
}

type Account struct {
	Username string `json:"username" xml:"username" yaml:"username" toml:"username"`
	Password string `json:"password" xml:"password" yaml:"password" toml:"password"`
	Dir      string `json:"dir" xml:"dir" yaml:"dir" toml:"dir"`
}

func (p *Config) init() error {
	if len(p.ServerName) == 0 {
		p.ServerName = "FTP Server"
	}
	if len(p.Dir) == 0 {
		p.Dir = "/"
	}
	if err := os.Mkdir(p.Dir, 0777); err != nil {
		return err
	}

	if p.Port == 0 {
		p.Port = 21
	}
	if p.DeadlineSeconds == 0 {
		p.DeadlineSeconds = 30
	}

	p.addr = GetAddress(p.Host, int(p.Port))

	if p.PasvMaxPort < p.PasvMinPort || p.PasvMaxPort > 65534 {
		return errors.New("params invalid, please check  pasv port")
	}

	p.accountMap = make(map[string]*Account, len(p.Accounts))
	if p.Accounts != nil && len(p.Accounts) > 0 {
		for _, account := range p.Accounts {
			if len(account.Username) > 0 && len(account.Password) > 0 {
				if account.Dir == ".." {
					account.Dir = ""
				}
				p.accountMap[account.Username] = account
			}
		}
	} else {
		p.accountMap["admin"] = &Account{"admin", "admin", ""}
	}

	return nil
}
