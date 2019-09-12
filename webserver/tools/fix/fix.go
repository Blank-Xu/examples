package fix

import (
	"time"
)

type Fix struct {
	TimeZone TimeZone `yaml:"time_zone"`
}

func (p *Fix) Init() {
	p.TimeZone.init()
}

type TimeZone struct {
	Name   string `json:"name" yaml:"name"`
	Offset int    `json:"offset" yaml:"offset"`
}

func (p *TimeZone) init() {
	if p == nil {
		return
	}

	if len(p.Name) == 0 {
		p.Name = "UTC"
	}
	if p.Offset <= 0 {
		p.Offset = 8
	}

	time.Local = time.FixedZone(p.Name, p.Offset*3600)
}
