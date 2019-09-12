package db

type Config struct {
	DriverName      string `json:"driver_name" yaml:"driver_name"`
	DataBase        string `json:"data_base" yaml:"data_base"`
	Host            string `json:"host" yaml:"host"`
	Port            string `json:"port" yaml:"port"`
	Username        string `json:"username" yaml:"username"`
	Password        string `json:"password" yaml:"password"`
	Charset         string `json:"charset" yaml:"charset"`
	LogLevel        int    `json:"log_level" yaml:"log_level"`
	ConnMaxLifetime int    `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	MaxIdleConns    int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int    `json:"max_open_conns" yaml:"max_open_conns"`
	ShowSql         bool   `json:"show_sql" yaml:"show_sql"`
	ShowExecTime    bool   `json:"show_exec_time" yaml:"show_exec_time"`
	Connect         bool   `json:"connect" yaml:"connect"`
}

func LoadConfig(filename string) error {

	return nil
}

func (p *Config) Init() {

}
