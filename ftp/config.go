package ftp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strings"
)

const version = "0.1.0"

type Config struct {
	ServerName             string     `json:"server_name" xml:"server_name" yaml:"server_name" toml:"server_name"`
	Dir                    string     `json:"dir" xml:"dir" yaml:"dir" toml:"dir"`
	Host                   string     `json:"host" xml:"host" yaml:"host" toml:"host"`
	Port                   uint32     `json:"port" xml:"port" yaml:"port" toml:"port"`
	MaxConnections         uint32     `json:"max_connections" xml:"max_connections" yaml:"max_connections" toml:"max_connections"`
	DeadlineSeconds        uint32     `json:"deadline_seconds" xml:"deadline_seconds" yaml:"deadline_seconds" toml:"deadline_seconds"`
	ReadDeadlineSeconds    uint32     `json:"read_deadline_seconds" xml:"read_deadline_seconds" yaml:"read_deadline_seconds" toml:"read_deadline_seconds"`
	WriteDeadlineSeconds   uint32     `json:"write_deadline_seconds" xml:"write_deadline_seconds" yaml:"write_deadline_seconds" toml:"write_deadline_seconds"`
	KeepAlivePeriodSeconds uint32     `json:"keep_alive_period_seconds" xml:"keep_alive_period_seconds" yaml:"keep_alive_period_seconds" toml:"keep_alive_period_seconds"`
	PasvMinPort            uint32     `json:"pasv_min_port" xml:"pasv_min_port" yaml:"pasv_min_port" toml:"pasv_min_port"`
	PasvMaxPort            uint32     `json:"pasv_max_port" xml:"pasv_max_port" yaml:"pasv_max_port" toml:"pasv_max_port"`
	AutoTls                bool       `json:"auto_tls" xml:"auto_tls" yaml:"auto_tls" toml:"auto_tls"`
	TlsKey                 string     `json:"tls_key" xml:"tls_key" yaml:"tls_key" toml:"tls_key"`
	Accounts               []*Account `json:"accounts" xml:"accounts" yaml:"accounts" toml:"accounts"`

	externalIP string
	addr       string
	tlscfg     *tls.Config
	accountMap map[string]*Account
}

type Account struct {
	Username string `json:"username" xml:"username" yaml:"username" toml:"username"`
	Password string `json:"password" xml:"password" yaml:"password" toml:"password"`
	Dir      string `json:"dir" xml:"dir" yaml:"dir" toml:"dir"`
}

func (p *Config) Check() (err error) {
	if len(p.ServerName) == 0 {
		p.ServerName = "FTP Server"
	}
	if len(p.Dir) == 0 {
		if p.Dir, err = os.Getwd(); err != nil {
			return err
		}
	} else if err = os.Mkdir(p.Dir, 0766); err != nil && !os.IsExist(err) {
		return
	}

	// if len(p.Host) == 0 {
	// 	p.Host = GetLocalIp()
	// }
	if p.Port == 0 {
		p.Port = 21
	}
	if p.DeadlineSeconds == 0 {
		p.DeadlineSeconds = 900
	}

	if p.externalIP, err = GetExternalIP(); err != nil {
		return err
	}
	p.externalIP = strings.ReplaceAll(p.externalIP, ",", ".")

	p.addr = GetAddress(p.Host, int(p.Port))

	if p.PasvMaxPort < p.PasvMinPort || p.PasvMaxPort > 65534 {
		return errors.New("params invalid, please check pasv port")
	}

	if err = p.checkTLS(); err != nil {
		return
	}
	if err = p.checkAccount(); err != nil {
		return
	}

	return
}

func (p *Config) checkTLS() error {
	if p.AutoTls {

	}
	return nil
}

func (p *Config) checkAccount() error {
	p.accountMap = make(map[string]*Account, len(p.Accounts))
	if p.Accounts == nil || len(p.Accounts) == 0 {
		p.accountMap["admin"] = &Account{"admin", "admin", ""}
		return nil
	}

	for _, account := range p.Accounts {
		if len(account.Username) > 0 && len(account.Password) > 0 {
			if strings.Contains(account.Dir, "..") {
				return fmt.Errorf("params invalid, please check account[%s] dir[%s] set, it must under the work dir[%s]",
					account.Username, account.Dir, p.Dir)
			}

			path := GetAbsPath(p.Dir, account.Dir)
			if err := os.MkdirAll(path, 0766); err != nil && !os.IsExist(err) {
				return err
			}

			p.accountMap[account.Username] = account
		}
	}

	return nil
}
