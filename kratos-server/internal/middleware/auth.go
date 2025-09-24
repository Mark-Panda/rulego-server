package middleware

import (
	"context"
	"strings"

	pb "kratos-server/api/rulego/v1"
	"kratos-server/internal/conf"
	"kratos-server/internal/service"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// JWT认证中间件
func JWT(config *conf.RuleGo, userSvc *service.UserService) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 如果不需要认证，直接通过
			if !config.Jwt.RequireAuth {
				return handler(ctx, req)
			}

			var token string

			// 从transport中获取token
			if tr, ok := transport.FromServerContext(ctx); ok {
				switch tr := tr.(type) {
				case *http.Transport:
					// 从Header中获取token
					authHeader := tr.RequestHeader().Get("Authorization")
					if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
						token = strings.TrimPrefix(authHeader, "Bearer ")
					}
					// 也可以从query参数中获取apiKey
					if token == "" {
						token = tr.Request().URL.Query().Get("apiKey")
					}
				case *grpc.Transport:
					// 从gRPC metadata中获取token
					if authHeader := tr.RequestHeader().Get("authorization"); authHeader != "" {
						if strings.HasPrefix(authHeader, "Bearer ") {
							token = strings.TrimPrefix(authHeader, "Bearer ")
						}
					}
				}
			}

			if token == "" {
				return nil, errors.Unauthorized("UNAUTHORIZED", "missing token")
			}

			// 验证token
			resp, err := userSvc.VerifyToken(ctx, &pb.VerifyTokenRequest{
				Token: token,
			})
			if err != nil || !resp.Valid {
				return nil, errors.Unauthorized("UNAUTHORIZED", "invalid token")
			}

			// 将用户信息添加到上下文
			ctx = context.WithValue(ctx, "username", resp.Username)

			return handler(ctx, req)
		}
	}
}

// CORS中间件
func CORS() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 添加CORS headers
			if tr, ok := transport.FromServerContext(ctx); ok {
				if httpTr, ok := tr.(*http.Transport); ok {
					w := httpTr.ReplyHeader()
					w.Set("Access-Control-Allow-Origin", "*")
					w.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
					w.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				}
			}
			return handler(ctx, req)
		}
	}
}

// 从上下文中获取用户名
func GetUsernameFromContext(ctx context.Context) string {
	if username, ok := ctx.Value("username").(string); ok {
		return username
	}
	return ""
}
