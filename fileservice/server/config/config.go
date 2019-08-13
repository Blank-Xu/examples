package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	Default *config

	configFile = flag.String("file", "config.yaml", "config file")
)

type config struct {
	RunMode string `yaml:"run_mode"`

	IP      string `yaml:"ip"`
	Port    int    `yaml:"port"`
	WorkDir string `yaml:"work_dir"`

	UploadLimit     int   `yaml:"upload_limit"`
	UploadMaxSize   int64 `yaml:"upload_max_size"`
	UploadChunkSize int64 `yaml:"upload_chunk_size"`

	DownloadLimit     int   `yaml:"download_limit"`
	DownloadChunkSize int64 `yaml:"download_chunk_size"`
}

func Init() {
	flag.Parse()

	log.Printf("open config file: %s\n", *configFile)

	file, err := os.OpenFile(*configFile, os.O_RDONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("open config file failed, err: %v", err))
	}

	var cfg config
	if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
		panic(fmt.Sprintf("decode config file failed, err: %v", err))
	}

	log.Printf("decode config success, config: %+v\n", cfg)

	Default = &cfg

	defaultCheck()

	log.Println("load config success")
}

func defaultCheck() {
	if len(Default.RunMode) == 0 {
		Default.RunMode = DEV
	}

	if Default.Port == 0 {
		Default.Port = 8080
	}

	if len(Default.WorkDir) == 0 {
		Default.WorkDir = "files"
	}

	// if !strings.ContainsRune(Default.WorkDir, os.PathSeparator) {
	// 	Default.WorkDir = fmt.Sprintf("%s%v", Default.WorkDir, os.PathSeparator)
	// }

	if err := os.Mkdir(Default.WorkDir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			panic(fmt.Sprintf("mkdir[%s] failed, err: %v", Default.WorkDir, err))
		}
	}

	if Default.UploadMaxSize == 0 {
		Default.UploadMaxSize = defaultMaxSize
	} else {
		Default.UploadMaxSize *= 1024 * 1024
	}

	if Default.UploadChunkSize == 0 {
		Default.UploadChunkSize = defaultChunkSize
	} else {
		Default.UploadChunkSize *= 1024 * 1024
	}

	// if Default.DownloadMaxSize == 0 {
	// 	Default.DownloadMaxSize = defaultMaxSize
	// } else{
	// 		Default.DownloadMaxSize *= 1024*1024
	// 	}

	if Default.DownloadChunkSize == 0 {
		Default.DownloadChunkSize = defaultChunkSize
	} else {
		Default.DownloadChunkSize *= 1024 * 1024
	}
}
