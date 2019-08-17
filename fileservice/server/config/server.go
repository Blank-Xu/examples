package config

type server struct {
	Name               string `yaml:"name"`
	IP                 string `yaml:"ip"`
	Port               int    `yaml:"port"`
	ReadTimeout        int    `yaml:"read_timeout"`
	WriteTimeout       int    `yaml:"write_timeout"`
	IdleTimeout        int    `yaml:"idle_timeout"`
	MaxConnPerHost     int    `yaml:"max_conn_per_host"`
	MaxIdleConn        int    `yaml:"max_idle_conn"`
	MaxIdleConnPerHost int    `yaml:"max_idle_conn_per_host"`
}
