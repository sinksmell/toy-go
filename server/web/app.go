package web

import (
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/server/xecho"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sinksmell/toy-go/handler/json"
)

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.serveHTTP,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

// HTTP地址
func (eng *Engine) serveHTTP() error {
	server := xecho.StdConfig("http").Build()
	server.Use(
		middleware.CORSWithConfig(middleware.DefaultCORSConfig),
	)

	server.POST("/api/json/to/go", json.ConvertToGo)
	server.Static("/", "./webapp")

	return eng.Serve(server)
}
