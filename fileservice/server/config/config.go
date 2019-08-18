package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

var (
	Default *config

	configFile = flag.String("file", "config.yaml", "config file")
)

type config struct {
	Server     server     `yaml:"server"`
	FileConfig FileConfig `yaml:"file_config"`
	TimeZone   TimeZone   `yaml:"time_zone"`
	LogConfig  LogConfig  `yaml:"log"`
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

	log.Printf("decode config success, config: %+v", cfg)

	initTime(cfg.TimeZone)

	initLog(cfg.LogConfig)

	Default = &cfg

	defaultCheck()

	logrus.Info("load config success")
}

func defaultCheck() {
	if Default.Server.Port == 0 {
		Default.Server.Port = 8080
	}

	if len(Default.FileConfig.WorkDir) == 0 {
		Default.FileConfig.WorkDir = "files"
	}

	if err := os.Mkdir(Default.FileConfig.WorkDir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			panic(fmt.Sprintf("mkdir[%s] failed, err: %v", Default.FileConfig.WorkDir, err))
		}
	}

	if Default.FileConfig.UploadMaxSize == 0 {
		Default.FileConfig.UploadMaxSize = defaultMaxSize
	} else {
		Default.FileConfig.UploadMaxSize *= 1024 * 1024
	}

	if Default.FileConfig.UploadChunkSize == 0 {
		Default.FileConfig.UploadChunkSize = defaultChunkSize
	} else {
		Default.FileConfig.UploadChunkSize *= 1024 * 1024
	}
}
