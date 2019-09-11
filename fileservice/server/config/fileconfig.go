package config

import (
	"fmt"
	"os"

	"fileservice/server/utils"
)

type FileConfig struct {
	WorkDir         string `yaml:"work_dir"`
	UploadLimit     int    `yaml:"upload_limit"`
	UploadMaxSize   int64  `yaml:"upload_max_size"`
	UploadChunkSize int64  `yaml:"upload_chunk_size"`
	DownloadLimit   int    `yaml:"download_limit"`
	FileMd5Limit    int    `yaml:"file_md5_limit"`
}

func (p *FileConfig) init() {
	if len(p.WorkDir) == 0 {
		p.WorkDir = "files"
	}

	if err := utils.MkdirAll(p.WorkDir); err != nil {
		if !os.IsExist(err) {
			panic(fmt.Sprintf("mkdir[%s] failed, err: %v", p.WorkDir, err))
		}
	}

	if p.UploadMaxSize == 0 {
		p.UploadMaxSize = defaultMaxSize
	} else {
		p.UploadMaxSize *= 1024 * 1024
	}

	if p.UploadChunkSize == 0 {
		p.UploadChunkSize = defaultChunkSize
	} else {
		p.UploadChunkSize *= 1024 * 1024
	}
}
