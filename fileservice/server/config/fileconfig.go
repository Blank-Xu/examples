package config

type FileConfig struct {
	WorkDir         string `yaml:"work_dir"`
	UploadLimit     int    `yaml:"upload_limit"`
	UploadMaxSize   int64  `yaml:"upload_max_size"`
	UploadChunkSize int64  `yaml:"upload_chunk_size"`
	DownloadLimit   int    `yaml:"download_limit"`
	FileMd5Limit    int    `yaml:"file_md5_limit"`
}
