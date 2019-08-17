package config

import (
	"time"
)

type TimeZone struct {
	Name   string `yaml:"name"`
	Offset int    `yaml:"offset"`
}

// 设置时区
func initTime(cfg TimeZone) {
	time.FixedZone(cfg.Name, cfg.Offset*3600)
}
