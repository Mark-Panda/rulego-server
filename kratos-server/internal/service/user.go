package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	pb "kratos-server/api/rulego/v1"
	"kratos-server/internal/biz"
	"kratos-server/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt"
)

// UserService 用户服务
type UserService struct {
	pb.UnimplementedUserServiceServer

	uc     *biz.UserUsecase
	config *conf.RuleGo
	log    *log.Helper
}

// NewUserService 创建用户服务
func NewUserService(uc *biz.UserUsecase, config *conf.Bootstrap, logger log.Logger) *UserService {
	return &UserService{
		uc:     uc,
		config: config.GetRulego(),
		log:    log.NewHelper(logger),
	}
}

// JWTClaim JWT声明
type JWTClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.uc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	// 生成JWT Token
	token, expiresAt, err := s.generateToken(user.Username)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token:     token,
		Username:  user.Username,
		ExpiresAt: expiresAt,
	}, nil
}

// VerifyToken 验证Token
func (s *UserService) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	claims, err := s.parseToken(req.Token)
	if err != nil {
		return &pb.VerifyTokenResponse{
			Valid: false,
		}, nil
	}

	return &pb.VerifyTokenResponse{
		Valid:     true,
		Username:  claims.Username,
		ExpiresAt: claims.StandardClaims.ExpiresAt,
	}, nil
}

// generateToken 生成JWT Token
func (s *UserService) generateToken(username string) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(s.config.Jwt.ExpireTime) * time.Millisecond).Unix()

	claims := &JWTClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  now.Unix(),
			Issuer:    s.config.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Jwt.SecretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

// parseToken 解析JWT Token
func (s *UserService) parseToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Jwt.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// hashPassword 哈希密码
func (s *UserService) hashPassword(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}
