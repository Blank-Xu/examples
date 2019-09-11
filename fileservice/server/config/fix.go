package config

import (
	"time"
)

type Fix struct {
	TimeZone TimeZone `yaml:"time_zone"`
}

func (p *Fix) init() {
	p.TimeZone.init()
}

type TimeZone struct {
	Name   string `yaml:"name"`
	Offset int    `yaml:"offset"`
}

// 设置时区
func (p *TimeZone) init() {
	if len(p.Name) == 0 {
		p.Name = "UTC"
	}
	if p.Offset <= 0 {
		p.Offset = 8
	}

	time.FixedZone(p.Name, p.Offset*3600)
}
