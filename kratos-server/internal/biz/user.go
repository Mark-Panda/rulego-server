package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// User 用户业务模型
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // 不在JSON中暴露密码
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepo 用户仓库接口
type UserRepo interface {
	FindByUsername(context.Context, string) (*User, error)
	FindByAPIKey(context.Context, string) (*User, error)
	Save(context.Context, *User) (*User, error)
	CheckPassword(context.Context, string, string) bool
}

// UserUsecase 用户用例
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// Login 用户登录
func (uc *UserUsecase) Login(ctx context.Context, username, password string) (*User, error) {
	uc.log.WithContext(ctx).Infof("Login: %v", username)

	if !uc.repo.CheckPassword(ctx, username, password) {
		return nil, ErrUnauthorized
	}

	return uc.repo.FindByUsername(ctx, username)
}

// VerifyToken 验证API Token
func (uc *UserUsecase) VerifyToken(ctx context.Context, apiKey string) (*User, error) {
	uc.log.WithContext(ctx).Infof("VerifyToken")
	return uc.repo.FindByAPIKey(ctx, apiKey)
}
