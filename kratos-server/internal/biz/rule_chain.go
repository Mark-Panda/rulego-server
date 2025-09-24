package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 定义业务层错误
var (
	ErrRuleNotFound      = errors.New("rule not found")
	ErrComponentNotFound = errors.New("component not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidParams     = errors.New("invalid parameters")
)

// RuleChain 规则链业务模型
type RuleChain struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Definition    string    `json:"definition"`
	Description   string    `json:"description"`
	Configuration string    `json:"configuration"` // 新增配置字段
	Status        int       `json:"status"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RuleChainRepo 规则链仓库接口
type RuleChainRepo interface {
	Save(context.Context, *RuleChain) (*RuleChain, error)
	Update(context.Context, *RuleChain) (*RuleChain, error)
	FindByID(context.Context, string, string) (*RuleChain, error)
	ListByUser(context.Context, string, int, int) ([]*RuleChain, int64, error)
	Delete(context.Context, string, string) error
}

// RuleChainUsecase 规则链用例
type RuleChainUsecase struct {
	repo      RuleChainRepo
	rulegoSvc RuleGoService // 注入RuleGo服务
	log       *log.Helper
}

// RuleGoService RuleGo服务接口
type RuleGoService interface {
	ExecuteRule(ctx context.Context, ruleID, username, msgType, data string, metadata map[string]string) (string, error)
	PostMessage(ctx context.Context, ruleID, username, msgType, data string, metadata map[string]string) error
	DeployRule(ctx context.Context, ruleID, username, ruleDefinition string) error
	UndeployRule(ctx context.Context, ruleID, username string) error
}

// NewRuleChainUsecase 创建规则链用例
func NewRuleChainUsecase(repo RuleChainRepo, rulegoSvc RuleGoService, logger log.Logger) *RuleChainUsecase {
	return &RuleChainUsecase{repo: repo, rulegoSvc: rulegoSvc, log: log.NewHelper(logger)}
}

// CreateRule 创建规则链
func (uc *RuleChainUsecase) CreateRule(ctx context.Context, rule *RuleChain) (*RuleChain, error) {
	uc.log.WithContext(ctx).Infof("CreateRule: %v", rule.ID)
	return uc.repo.Save(ctx, rule)
}

// UpdateRule 更新规则链
func (uc *RuleChainUsecase) UpdateRule(ctx context.Context, rule *RuleChain) (*RuleChain, error) {
	uc.log.WithContext(ctx).Infof("UpdateRule: %v", rule.ID)
	return uc.repo.Update(ctx, rule)
}

// GetRule 获取规则链
func (uc *RuleChainUsecase) GetRule(ctx context.Context, id, username string) (*RuleChain, error) {
	uc.log.WithContext(ctx).Infof("GetRule: %v", id)
	return uc.repo.FindByID(ctx, id, username)
}

// ListRules 获取规则链列表
func (uc *RuleChainUsecase) ListRules(ctx context.Context, username string, page, pageSize int) ([]*RuleChain, int64, error) {
	uc.log.WithContext(ctx).Infof("ListRules: %v", username)
	return uc.repo.ListByUser(ctx, username, page, pageSize)
}

// DeleteRule 删除规则链
func (uc *RuleChainUsecase) DeleteRule(ctx context.Context, id, username string) error {
	uc.log.WithContext(ctx).Infof("DeleteRule: %v", id)
	return uc.repo.Delete(ctx, id, username)
}

// ExecuteRule 执行规则链
func (uc *RuleChainUsecase) ExecuteRule(ctx context.Context, id, username, msgType, data string, metadata map[string]string) (string, error) {
	uc.log.WithContext(ctx).Infof("ExecuteRule: %v", id)

	// 首先检查规则链是否存在
	rule, err := uc.repo.FindByID(ctx, id, username)
	if err != nil {
		return "", fmt.Errorf("rule not found: %w", err)
	}

	// 调用RuleGo服务执行规则链
	result, err := uc.rulegoSvc.ExecuteRule(ctx, id, username, msgType, data, metadata)
	if err != nil {
		// 记录执行失败日志
		uc.log.WithContext(ctx).Errorf("Failed to execute rule %s: %v", id, err)
		return "", fmt.Errorf("rule execution failed: %w", err)
	}

	uc.log.WithContext(ctx).Infof("Rule %s executed successfully for user %s", rule.Name, username)
	return result, nil
}

// PostMessage 异步推送消息到规则链
func (uc *RuleChainUsecase) PostMessage(ctx context.Context, id, username, msgType, data string, metadata map[string]string) error {
	uc.log.WithContext(ctx).Infof("PostMessage: %v", id)

	// 首先检查规则链是否存在
	rule, err := uc.repo.FindByID(ctx, id, username)
	if err != nil {
		return fmt.Errorf("rule not found: %w", err)
	}

	// 调用RuleGo服务异步推送消息
	err = uc.rulegoSvc.PostMessage(ctx, id, username, msgType, data, metadata)
	if err != nil {
		// 记录推送失败日志
		uc.log.WithContext(ctx).Errorf("Failed to post message to rule %s: %v", id, err)
		return fmt.Errorf("message posting failed: %w", err)
	}

	uc.log.WithContext(ctx).Infof("Message posted successfully to rule %s for user %s", rule.Name, username)
	return nil
}

// DeployRule 部署规则链
func (uc *RuleChainUsecase) DeployRule(ctx context.Context, id, username string) error {
	uc.log.WithContext(ctx).Infof("DeployRule: %v", id)

	// 首先获取规则链定义
	rule, err := uc.repo.FindByID(ctx, id, username)
	if err != nil {
		return fmt.Errorf("rule not found: %w", err)
	}

	// 检查规则链定义是否为空
	if rule.Definition == "" {
		return fmt.Errorf("rule definition is empty for rule %s", id)
	}

	// 调用RuleGo服务部署规则链
	err = uc.rulegoSvc.DeployRule(ctx, id, username, rule.Definition)
	if err != nil {
		// 记录部署失败日志
		uc.log.WithContext(ctx).Errorf("Failed to deploy rule %s: %v", id, err)
		return fmt.Errorf("rule deployment failed: %w", err)
	}

	// 更新规则链状态为已部署（1）
	rule.Status = 1
	_, err = uc.repo.Update(ctx, rule)
	if err != nil {
		uc.log.WithContext(ctx).Warnf("Failed to update rule status after deployment: %v", err)
		// 不返回错误，因为部署已经成功
	}

	uc.log.WithContext(ctx).Infof("Rule %s deployed successfully for user %s", rule.Name, username)
	return nil
}

// UndeployRule 下线规则链
func (uc *RuleChainUsecase) UndeployRule(ctx context.Context, id, username string) error {
	uc.log.WithContext(ctx).Infof("UndeployRule: %v", id)

	// 首先检查规则链是否存在
	rule, err := uc.repo.FindByID(ctx, id, username)
	if err != nil {
		return fmt.Errorf("rule not found: %w", err)
	}

	// 调用RuleGo服务下线规则链
	err = uc.rulegoSvc.UndeployRule(ctx, id, username)
	if err != nil {
		// 记录下线失败日志
		uc.log.WithContext(ctx).Errorf("Failed to undeploy rule %s: %v", id, err)
		return fmt.Errorf("rule undeployment failed: %w", err)
	}

	// 更新规则链状态为已停用（0）
	rule.Status = 0
	_, err = uc.repo.Update(ctx, rule)
	if err != nil {
		uc.log.WithContext(ctx).Warnf("Failed to update rule status after undeployment: %v", err)
		// 不返回错误，因为下线已经成功
	}

	uc.log.WithContext(ctx).Infof("Rule %s undeployed successfully for user %s", rule.Name, username)
	return nil
}
