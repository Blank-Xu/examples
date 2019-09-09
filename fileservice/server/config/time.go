package config

import (
	"time"
)

type TimeZone struct {
	Name   string `yaml:"name"`
	Offset int    `yaml:"offset"`
}

// 设置时区
func (p *TimeZone) init() {
	time.FixedZone(p.Name, p.Offset*3600)
}
