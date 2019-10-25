package ftp

import (
	"os"
	"time"
)

func NewDirInfo(name string, modTime time.Time) os.FileInfo {
	return &FileInfo{name: name, mode: os.ModeDir | 0666, modTime: modTime}
}

func NewFileInfo(name string, size int64, modTime time.Time) os.FileInfo {
	return &FileInfo{name: name, size: size, mode: 0666, modTime: modTime}
}

type FileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (p *FileInfo) Name() string {
	return p.name
}

func (p *FileInfo) Size() int64 {
	return p.size
}

func (p *FileInfo) Mode() os.FileMode {
	return p.mode
}

func (p *FileInfo) ModTime() time.Time {
	return p.modTime
}

func (p *FileInfo) IsDir() bool {
	return (p.mode & os.ModeDir) == os.ModeDir
}

func (p *FileInfo) Sys() interface{} {
	return nil
}
