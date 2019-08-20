package config

type server struct {
	Name                string `yaml:"name"`
	IP                  string `yaml:"ip"`
	Port                int    `yaml:"port"`
	ReadTimeout         int    `yaml:"read_timeout"`
	WriteTimeout        int    `yaml:"write_timeout"`
	IdleTimeout         int    `yaml:"idle_timeout"`
	MaxConnsPerHost     int    `yaml:"max_conns_per_host"`      // 每一个host对应的最大连接数
	MaxIdleConns        int    `yaml:"max_idle_conns"`          // 所有host对应的idle状态最大的连接总数
	MaxIdleConnsPerHost int    `yaml:"max_idle_conns_per_host"` // 每一个host对应idle状态的最大的连接数
}
