package server

import (
	v1 "kratos-server/api/rulego/v1"
	"kratos-server/internal/conf"
	"kratos-server/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *conf.Server,
	userSvc *service.UserService,
	ruleSvc *service.RuleService,
	componentSvc *service.ComponentService,
	logSvc *service.LogService,
	localeSvc *service.LocaleService,
	logger log.Logger,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	// 注册HTTP服务
	v1.RegisterUserServiceHTTPServer(srv, userSvc)
	v1.RegisterRuleServiceHTTPServer(srv, ruleSvc)
	v1.RegisterComponentServiceHTTPServer(srv, componentSvc)
	v1.RegisterLogServiceHTTPServer(srv, logSvc)
	v1.RegisterLocaleServiceHTTPServer(srv, localeSvc)

	return srv
}
