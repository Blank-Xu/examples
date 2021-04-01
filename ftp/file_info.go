package ftp

import (
	"os"
	"time"
)

// NewDirInfo .
func NewDirInfo(name string, modTime time.Time) os.FileInfo {
	return &FileInfo{name: name, mode: os.ModeDir | 0666, modTime: modTime}
}

// NewFileInfo .
func NewFileInfo(name string, size int64, modTime time.Time) os.FileInfo {
	return &FileInfo{name: name, size: size, mode: 0666, modTime: modTime}
}

// FileInfo .
type FileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name .
func (p *FileInfo) Name() string {
	return p.name
}

// Size .
func (p *FileInfo) Size() int64 {
	return p.size
}

// Mode .
func (p *FileInfo) Mode() os.FileMode {
	return p.mode
}

// ModTime .
func (p *FileInfo) ModTime() time.Time {
	return p.modTime
}

// IsDir .
func (p *FileInfo) IsDir() bool {
	return (p.mode & os.ModeDir) == os.ModeDir
}

func (p *FileInfo) Sys() interface{} {
	return nil
}
