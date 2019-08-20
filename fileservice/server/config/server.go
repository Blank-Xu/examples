package config

type server struct {
	Name               string `yaml:"name"`
	IP                 string `yaml:"ip"`
	Port               int    `yaml:"port"`
	ReadTimeout        int    `yaml:"read_timeout"`
	WriteTimeout       int    `yaml:"write_timeout"`
	IdleTimeout        int    `yaml:"idle_timeout"`
	MaxConnPerHost     int    `yaml:"max_conn_per_host"`      // 每一个host对应的最大连接数
	MaxIdleConn        int    `yaml:"max_idle_conn"`          // 所有hosts对应的最大的连接总数
	MaxIdleConnPerHost int    `yaml:"max_idle_conn_per_host"` // 每一个host对应的最大的空闲连接数
}
