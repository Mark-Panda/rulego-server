package service

import (
	"context"
	"kratos-server/internal/biz"
)

// RuleGoAdapter 实现biz层的RuleGoService接口
type RuleGoAdapter struct {
	rulegoSvc *RuleGoService
}

// NewRuleGoAdapter 创建RuleGo适配器
func NewRuleGoAdapter(rulegoSvc *RuleGoService) biz.RuleGoService {
	return &RuleGoAdapter{
		rulegoSvc: rulegoSvc,
	}
}

// ExecuteRule 执行规则链
func (a *RuleGoAdapter) ExecuteRule(ctx context.Context, ruleID, username, msgType, data string, metadata map[string]string) (string, error) {
	return a.rulegoSvc.ExecuteRule(ctx, ruleID, username, msgType, data, metadata)
}

// PostMessage 异步推送消息到规则链
func (a *RuleGoAdapter) PostMessage(ctx context.Context, ruleID, username, msgType, data string, metadata map[string]string) error {
	return a.rulegoSvc.PostMessage(ctx, ruleID, username, msgType, data, metadata)
}

// DeployRule 部署规则链
func (a *RuleGoAdapter) DeployRule(ctx context.Context, ruleID, username, ruleDefinition string) error {
	return a.rulegoSvc.DeployRule(ctx, ruleID, username, ruleDefinition)
}

// UndeployRule 下线规则链
func (a *RuleGoAdapter) UndeployRule(ctx context.Context, ruleID, username string) error {
	return a.rulegoSvc.UndeployRule(ctx, ruleID, username)
}
