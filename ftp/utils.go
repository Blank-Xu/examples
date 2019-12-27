package ftp

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	random = rand.New(rand.NewSource(time.Now().UnixNano()))

	defaultHttpClient = http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

func GetExternalIP() (string, error) {
	resp, err := defaultHttpClient.Get("http://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	buf.Grow(64)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func GetLocalIp() string {
	netInterfaces, _ := net.Interfaces()

	for i := 0; i < len(netInterfaces); i++ {
		if netInterfaces[i].Flags&net.FlagUp == 0 {
			continue
		}

		addrs, _ := netInterfaces[i].Addrs()
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	return ""
}

func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}

	if ip4 := IP.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}

	return false
}

func GetAddress(host string, port int) (addr string) {
	var buf bytes.Buffer
	buf.Grow(len(host) * 2)

	if strings.IndexByte(host, ':') > -1 {
		buf.WriteByte('[')
		buf.WriteString(host)
		buf.WriteByte(']')
	} else {
		buf.WriteString(host)
	}

	if port > 0 {
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(port))
	}

	return buf.String()
}

func GetRandomPort(minPort, maxPort int) (port int) {
	n := maxPort - minPort
	if n == 0 {
		if minPort > 0 {
			port = minPort
		}
	} else if n > 0 {
		port = minPort + random.Intn(n+1)
	}
	return
}

func NewActiveTCPConn(host string, port int) (*net.TCPConn, error) {
	addr := GetAddress(host, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	return net.DialTCP("tcp", nil, tcpAddr)
}

func NewPassiveTCPListener(host string, minPort, maxPort int) (*net.TCPListener, error) {
	for i := 0; i < 100; i++ {
		listener, err := NewTcpListener(host, GetRandomPort(minPort, maxPort))
		if err == nil {
			return listener, nil
		}
	}
	return nil, errors.New("Unable to find an available port")
}

func NewTcpListener(host string, port int) (*net.TCPListener, error) {
	addr := GetAddress(host, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", tcpAddr)
}

func GetFileList(absPath, path string) (files []os.FileInfo, err error) {
	now := time.Now().UTC()

	switch path {
	case "/debug":
		return
	case "/virtual":
		files = []os.FileInfo{
			NewFileInfo("localpath.txt", 1024, now),
			NewFileInfo("file2.txt", 2048, now),
		}
	default:
		if files, err = ioutil.ReadDir(absPath); err != nil {
			return
		}
		if path == "/" {
			files = append(files, NewDirInfo("virtual", now))
		}
	}
	return
}

func GetAbsPath(workDir, path string) string {
	if workDir == "./" || workDir == "/" {
		workDir = ""
	}

	if len(path) == 0 || path == "/" || path == "." {
		return workDir
	}

	newPath := filepath.Join(workDir, path)
	l := len(workDir)
	if len(newPath) < l || newPath[:l] != workDir {
		return workDir
	}

	return newPath
}

const (
	dateFormatStatTime      = "Jan _2 15:04"          // LIST date formatting with hour and minute
	dateFormatStatYear      = "Jan _2  2006"          // LIST date formatting with year
	dateFormatStatOldSwitch = time.Hour * 24 * 30 * 6 // 6 months ago
	dateFormatMLSD          = "20060102150405"        // MLSD date formatting
)

func GetFileModTime(baseTime, modTime time.Time) string {
	if baseTime.Sub(modTime) > dateFormatStatOldSwitch {
		return modTime.Format(dateFormatStatYear)
	}

	return modTime.Format(dateFormatStatTime)
}
