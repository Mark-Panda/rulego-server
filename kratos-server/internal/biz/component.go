package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// Component 组件业务模型
type Component struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Author      string    `json:"author"`
	DSL         string    `json:"dsl"`
	IsInstalled bool      `json:"is_installed"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ComponentRepo 组件仓库接口
type ComponentRepo interface {
	Save(context.Context, *Component) (*Component, error)
	Update(context.Context, *Component) (*Component, error)
	FindByID(context.Context, string, string) (*Component, error)
	ListByUser(context.Context, string) ([]*Component, error)
	Delete(context.Context, string, string) error
}

// ComponentUsecase 组件用例
type ComponentUsecase struct {
	repo ComponentRepo
	log  *log.Helper
}

// NewComponentUsecase 创建组件用例
func NewComponentUsecase(repo ComponentRepo, logger log.Logger) *ComponentUsecase {
	return &ComponentUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateComponent 创建组件
func (uc *ComponentUsecase) CreateComponent(ctx context.Context, component *Component) (*Component, error) {
	uc.log.WithContext(ctx).Infof("CreateComponent: %v", component.ID)
	return uc.repo.Save(ctx, component)
}

// UpdateComponent 更新组件
func (uc *ComponentUsecase) UpdateComponent(ctx context.Context, component *Component) (*Component, error) {
	uc.log.WithContext(ctx).Infof("UpdateComponent: %v", component.ID)
	return uc.repo.Update(ctx, component)
}

// GetComponent 获取组件
func (uc *ComponentUsecase) GetComponent(ctx context.Context, id, username string) (*Component, error) {
	uc.log.WithContext(ctx).Infof("GetComponent: %v", id)
	return uc.repo.FindByID(ctx, id, username)
}

// ListComponents 获取组件列表
func (uc *ComponentUsecase) ListComponents(ctx context.Context, username string) ([]*Component, error) {
	uc.log.WithContext(ctx).Infof("ListComponents: %v", username)
	return uc.repo.ListByUser(ctx, username)
}

// DeleteComponent 删除组件
func (uc *ComponentUsecase) DeleteComponent(ctx context.Context, id, username string) error {
	uc.log.WithContext(ctx).Infof("DeleteComponent: %v", id)
	return uc.repo.Delete(ctx, id, username)
}

// RunLog 运行日志业务模型
type RunLog struct {
	ID           uint      `json:"id"`
	ChainID      string    `json:"chain_id"`
	MsgType      string    `json:"msg_type"`
	Data         string    `json:"data"`
	Result       string    `json:"result"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message"`
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
}

// RunLogRepo 运行日志仓库接口
type RunLogRepo interface {
	Save(context.Context, *RunLog) (*RunLog, error)
	ListByUser(context.Context, string, int, int) ([]*RunLog, int64, error)
	Delete(context.Context, []uint, string) (int64, error)
}

// RunLogUsecase 运行日志用例
type RunLogUsecase struct {
	repo RunLogRepo
	log  *log.Helper
}

// NewRunLogUsecase 创建运行日志用例
func NewRunLogUsecase(repo RunLogRepo, logger log.Logger) *RunLogUsecase {
	return &RunLogUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateRunLog 创建运行日志
func (uc *RunLogUsecase) CreateRunLog(ctx context.Context, runLog *RunLog) (*RunLog, error) {
	uc.log.WithContext(ctx).Infof("CreateRunLog: %v", runLog.ChainID)
	return uc.repo.Save(ctx, runLog)
}

// ListRunLogs 获取运行日志列表
func (uc *RunLogUsecase) ListRunLogs(ctx context.Context, username string, page, pageSize int) ([]*RunLog, int64, error) {
	uc.log.WithContext(ctx).Infof("ListRunLogs: %v", username)
	return uc.repo.ListByUser(ctx, username, page, pageSize)
}

// DeleteRunLogs 删除运行日志
func (uc *RunLogUsecase) DeleteRunLogs(ctx context.Context, ids []uint, username string) (int64, error) {
	uc.log.WithContext(ctx).Infof("DeleteRunLogs: %v", ids)
	return uc.repo.Delete(ctx, ids, username)
}
