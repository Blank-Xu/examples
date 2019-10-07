package config

type config struct {
	WindowHeight int    `yaml:"window_height"`
	WindowWidth  int    `yaml:"window_width"`
	DownloadFile string `yaml:"download_file"`
	ThreadCount  int    `yaml:"thread_count"`
	RetryCount   int    `yaml:"retry_count"`
	Proxy        string `yaml:"proxy"`
	Socks5Proxy  string `yaml:"socks5_proxy"`
	Stream       string `yaml:"stream"`
	ChunkSize    int    `yaml:"chunk_size"`
	// UseAria2RPC Use Aria2 RPC to download
	UseAria2RPC bool
	// Aria2Token Aria2 RPC Token
	Aria2Token string
	// Aria2Addr Aria2 Address (default "localhost:6800")
	Aria2Addr string
	// Aria2Method Aria2 Method (default "http")
	Aria2Method string

	YouKu *Youku `yaml:"you_ku"`
}
