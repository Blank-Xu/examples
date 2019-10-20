package ftp

import (
	"time"
)

type Logger struct {
	User          string    `json:"user"`
	Addr          string    `json:"addr"`
	Command       string    `json:"command"`
	Path          string    `json:"path"`
	ConnectedTime time.Time `json:"connected_time"`
	Time          time.Time `json:"time"`
}

var DefaultLogger *Logger

func (p *Logger) log(level int, msg string) {

}

func (p *Logger) Debug(msg string) {
	p.log(0, msg)
}
