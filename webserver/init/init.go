package init

import (
	"flag"
	"io"
	"log"
	"os"
	"time"

	"webserver/tools/fix"
	"webserver/tools/httpserver"
	"webserver/tools/jwt"
	wlog "webserver/tools/log"
)

const (
	ConfigDir  = "configs"
	ServerName = "webserver"
	Version    = "0.1.0"
)

const (
	remoteConfigRefreshTime = time.Second * 5
)

var (
	runMode = flag.String("m", "debug", "run mode")

	// for remote config setting
	remote        = flag.Bool("r", false, "remote")
	provider      = flag.String("P", "etcd", "provider")
	endpoint      = flag.String("h", "http://127.0.0.1:4001", "endpoint")
	wpath         = flag.String("p", ServerName, "path")
	secretkeyring = flag.String("s", ServerName, "secretkeyring")
	configType    = flag.String("t", "yaml", "configType")
)

type config struct {
	Fix    *fix.Fix               `json:"fix" yaml:"fix"`
	Server *httpserver.HttpServer `json:"server" yaml:"server"`
	Jwt    []*jwt.Jwt             `json:"jwt" yaml:"jwt"`
	Log    *wlog.Log              `json:"log" yaml:"log"`
}

func Init() {
	flag.Parse()
	log.Printf("server version: %s, start args: %v\n", Version, flag.Args())

	if *remote {
		parseRemoteConfig(loadConfig)
	} else {
		parseLocalConfig(loadConfig)
	}

	log.SetOutput(io.MultiWriter(os.Stdout))
}
