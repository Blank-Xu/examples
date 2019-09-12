package log

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"

	"webserver/utils"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

var Default *zap.Logger

type Log struct {
	WorkDir      string `json:"work_dir" yaml:"work_dir"`
	WriteFile    bool   `json:"write_file" yaml:"write_file"`
	Filename     string `json:"filename" yaml:"filename"`
	LogLevel     uint8  `json:"log_level" yaml:"log_level"`
	Linkname     string `json:"linkname" yaml:"linkname"`
	TimeFormat   string `json:"time_format" yaml:"time_format"`
	ReportCaller bool   `json:"report_caller" yaml:"report_caller"`

	MaxAge       int  `json:"max_age" yaml:"max_age"`             // 保存天数，单位：天
	RotationTime int  `json:"rotation_time" yaml:"rotation_time"` // 分割时间，单位：小时
	JsonFormat   bool `json:"json_format" yaml:"json_format"`
}

func (p *Log) Init() (err error) {
	Default, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}

func (p *Log) getWriter() io.Writer {
	if p.WriteFile {
		if len(p.WorkDir) > 0 {
			if err := utils.MkdirAll(p.WorkDir); err != nil {
				log.Printf("mkdir [%s] failed, err: %v", p.WorkDir, err)
				panic(err)
			}
		}

		var rotate, err = rotatelogs.New(
			filepath.Join(p.WorkDir, p.Filename),
			rotatelogs.WithLinkName(filepath.Join(p.WorkDir, p.Linkname)),
			rotatelogs.WithRotationTime(time.Hour*time.Duration(p.RotationTime)),
			rotatelogs.WithMaxAge(time.Hour*time.Duration(24*p.MaxAge)),
			rotatelogs.WithLocation(time.Local),
		)
		if err != nil {
			log.Printf("create rotate log failed, err: %v", err)
			panic(err)
		}
		return rotate
	}

	return os.Stdout
}
