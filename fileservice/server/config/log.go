package config

import (
	"framework/fileservice/server/utils"
	"log"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type Log struct {
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

func (p *Log) init() {
	var format logrus.Formatter
	if p.JsonFormat {
		format = &logrus.JSONFormatter{TimestampFormat: p.TimeFormat}
	} else {
		format = &logrus.TextFormatter{TimestampFormat: p.TimeFormat}
	}
	logrus.SetFormatter(format)

	logrus.SetReportCaller(p.ReportCaller)

	if p.WriteFile {
		if len(p.WorkDir) > 0 {
			if err := utils.MkdirAll(p.WorkDir); err != nil {
				log.Printf("mkdir [%s] failed, err: %v", p.WorkDir, err)
				panic(err)
			}
		}

		var filename = filepath.Join(p.WorkDir, p.Filename)
		rotate, err := rotatelogs.New(
			filename,
			rotatelogs.WithLinkName(filepath.Join(p.WorkDir, p.Linkname)),
			rotatelogs.WithRotationTime(time.Hour*time.Duration(p.RotationTime)),
			rotatelogs.WithMaxAge(time.Hour*time.Duration(24*p.MaxAge)),
			rotatelogs.WithLocation(time.Local),
		)
		if err != nil {
			log.Println("create rotate log failed, err:", err)
			panic(err)
		}

		logrus.SetOutput(rotate)
	}
}
