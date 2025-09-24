//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"kratos-server/internal/biz"
	"kratos-server/internal/conf"
	"kratos-server/internal/data"
	"kratos-server/internal/server"
	"kratos-server/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// 提供器函数用于从 Bootstrap 中提取配置
func newServerConf(bc *conf.Bootstrap) *conf.Server {
	return bc.Server
}

func newDataConf(bc *conf.Bootstrap) *conf.Data {
	return bc.Data
}

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		newServerConf,
		newDataConf,
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}
