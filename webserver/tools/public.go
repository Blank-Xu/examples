package tools

import (
	"github.com/go-xorm/xorm"
	"go.uber.org/zap"

	"webserver/tools/httpserver"
)

var (
	HttpServer *httpserver.HttpServer

	Logger, _ = zap.NewProduction()
	DB        *xorm.Engine
)
