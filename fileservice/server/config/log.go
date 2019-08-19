package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	WorkDir      string `yaml:"work_dir"`
	WriteFile    bool   `yaml:"write_file"`
	Filename     string `yaml:"filename"`
	Linkname     string `yaml:"linkname"`
	TimeFormat   string `yaml:"time_format"`
	ReportCaller bool   `yaml:"report_caller"`
	MaxAge       int    `yaml:"max_age"`       // 保存天数，单位：天
	RotationTime int    `yaml:"rotation_time"` // 分割时间，单位：小时
	JsonFormat   bool   `yaml:"json_format"`
}

func initLog(cfg LogConfig) {
	var format logrus.Formatter
	if cfg.JsonFormat {
		format = &logrus.JSONFormatter{TimestampFormat: cfg.TimeFormat}
	} else {
		format = &logrus.TextFormatter{TimestampFormat: cfg.TimeFormat}
	}
	logrus.SetFormatter(format)

	logrus.SetReportCaller(cfg.ReportCaller)

	if cfg.WriteFile {
		if len(cfg.WorkDir) > 0 {
			if err := os.MkdirAll(cfg.WorkDir, 0666); err != nil {
				log.Println("mkdir failed, err:", err)
				panic(err)
			}
		}

		var filename = filepath.Join(cfg.WorkDir, cfg.Filename)
		rotate, err := rotatelogs.New(
			filename,
			rotatelogs.WithLinkName(filepath.Join(cfg.WorkDir, cfg.Linkname)),
			rotatelogs.WithRotationTime(time.Hour*time.Duration(cfg.RotationTime)),
			rotatelogs.WithMaxAge(time.Hour*time.Duration(24*cfg.MaxAge)),
			rotatelogs.WithLocation(time.Local),
		)
		if err != nil {
			log.Println("create rotate log failed, err:", err)
			panic(err)
		}

		logrus.SetOutput(rotate)
	}
}
