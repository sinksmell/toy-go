package main

import (
	"github.com/sinksmell/toy-go/web"

	"github.com/douyu/jupiter/pkg/xlog"
)

func main() {
	if err := web.NewEngine().Run(); err != nil {
		xlog.Error(err.Error())
	}
}
