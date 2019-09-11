package config

import (
	"flag"
	"fmt"
	"io"
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
	Fix        *Fix        `yaml:"fix"`
	Server     *Server     `yaml:"server"`
	FileConfig *FileConfig `yaml:"file_config"`
	Jwt        *Jwt        `yaml:"jwt"`
	Log        *Log        `yaml:"log"`
}

func Init() {
	flag.Parse()

	log.Printf("open config file: %s\n", *configFile)

	file, err := os.OpenFile(*configFile, os.O_RDONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("open config file failed, err: %v", err))
	}
	defer file.Close()

	var cfg config
	if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
		panic(fmt.Sprintf("decode config file failed, err: %v", err))
	}

	cfg.Fix.init()
	cfg.Log.init()
	cfg.Server.init()
	cfg.FileConfig.init()
	cfg.Jwt.init()

	log.SetOutput(io.MultiWriter(os.Stdout, logrus.StandardLogger().Writer()))

	Default = &cfg

	logrus.Info("load config success")
}
