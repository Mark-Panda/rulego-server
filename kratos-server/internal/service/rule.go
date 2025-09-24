package service

import (
	"context"
	"fmt"
	"time"

	pb "kratos-server/api/rulego/v1"
	"kratos-server/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// RuleService 规则服务
type RuleService struct {
	pb.UnimplementedRuleServiceServer

	uc  *biz.RuleChainUsecase
	log *log.Helper
}

// NewRuleService 创建规则服务
func NewRuleService(uc *biz.RuleChainUsecase, logger log.Logger) *RuleService {
	return &RuleService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// ListRules 获取规则链列表
func (s *RuleService) ListRules(ctx context.Context, req *pb.ListRulesRequest) (*pb.ListRulesResponse, error) {
	rules, total, err := s.uc.ListRules(ctx, req.Username, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	pbRules := make([]*pb.RuleInfo, len(rules))
	for i, rule := range rules {
		pbRules[i] = &pb.RuleInfo{
			Id:          rule.ID,
			Name:        rule.Name,
			CreatedTime: rule.CreatedAt.Unix(),
			UpdatedTime: rule.UpdatedAt.Unix(),
			Description: rule.Description,
			Status:      int32(rule.Status),
		}
	}

	return &pb.ListRulesResponse{
		Rules:    pbRules,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetRule 获取规则链详情
func (s *RuleService) GetRule(ctx context.Context, req *pb.GetRuleRequest) (*pb.GetRuleResponse, error) {
	rule, err := s.uc.GetRule(ctx, req.Id, req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.GetRuleResponse{
		Rule: &pb.RuleDetail{
			Id:          rule.ID,
			Name:        rule.Name,
			Definition:  rule.Definition,
			CreatedTime: rule.CreatedAt.Unix(),
			UpdatedTime: rule.UpdatedAt.Unix(),
			Description: rule.Description,
			Status:      int32(rule.Status),
		},
	}, nil
}

// GetLatestRule 获取最新的规则链
func (s *RuleService) GetLatestRule(ctx context.Context, req *pb.GetLatestRuleRequest) (*pb.GetLatestRuleResponse, error) {
	rule, err := s.uc.GetRule(ctx, req.Id, req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.GetLatestRuleResponse{
		Rule: &pb.RuleDetail{
			Id:          rule.ID,
			Name:        rule.Name,
			Definition:  rule.Definition,
			CreatedTime: rule.CreatedAt.Unix(),
			UpdatedTime: rule.UpdatedAt.Unix(),
			Description: rule.Description,
			Status:      int32(rule.Status),
		},
	}, nil
}

// SaveRule 保存规则链
func (s *RuleService) SaveRule(ctx context.Context, req *pb.SaveRuleRequest) (*pb.SaveRuleResponse, error) {
	rule := &biz.RuleChain{
		ID:          req.Id,
		Name:        req.Name,
		Definition:  req.Definition,
		Description: req.Description,
		Username:    req.Username,
	}

	_, err := s.uc.CreateRule(ctx, rule)
	if err != nil {
		return &pb.SaveRuleResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SaveRuleResponse{
		Success: true,
		Message: "Rule saved successfully",
	}, nil
}

// DeleteRule 删除规则链
func (s *RuleService) DeleteRule(ctx context.Context, req *pb.DeleteRuleRequest) (*pb.DeleteRuleResponse, error) {
	err := s.uc.DeleteRule(ctx, req.Id, req.Username)
	if err != nil {
		return &pb.DeleteRuleResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteRuleResponse{
		Success: true,
		Message: "Rule deleted successfully",
	}, nil
}

// SaveRuleBaseInfo 保存规则链基本信息
func (s *RuleService) SaveRuleBaseInfo(ctx context.Context, req *pb.SaveRuleBaseInfoRequest) (*pb.SaveRuleBaseInfoResponse, error) {
	// 获取现有规则链
	existingRule, err := s.uc.GetRule(ctx, req.Id, req.Username)
	if err != nil {
		return &pb.SaveRuleBaseInfoResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 更新基本信息
	existingRule.Name = req.Name
	existingRule.Description = req.Description

	_, err = s.uc.UpdateRule(ctx, existingRule)
	if err != nil {
		return &pb.SaveRuleBaseInfoResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SaveRuleBaseInfoResponse{
		Success: true,
		Message: "Rule base info saved successfully",
	}, nil
}

// SaveRuleConfiguration 保存规则链配置
func (s *RuleService) SaveRuleConfiguration(ctx context.Context, req *pb.SaveRuleConfigurationRequest) (*pb.SaveRuleConfigurationResponse, error) {
	// 获取现有规则链
	existingRule, err := s.uc.GetRule(ctx, req.Id, req.Username)
	if err != nil {
		return &pb.SaveRuleConfigurationResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 更新配置信息 - 这里可以根据var_type处理不同类型的配置
	existingRule.Configuration = req.Configuration

	_, err = s.uc.UpdateRule(ctx, existingRule)
	if err != nil {
		return &pb.SaveRuleConfigurationResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SaveRuleConfigurationResponse{
		Success: true,
		Message: "Rule configuration saved successfully",
	}, nil
}

// ExecuteRule 执行规则链
func (s *RuleService) ExecuteRule(ctx context.Context, req *pb.ExecuteRuleRequest) (*pb.ExecuteRuleResponse, error) {
	// 通过RuleGo执行规则链
	result, err := s.uc.ExecuteRule(ctx, req.Id, req.Username, req.MsgType, req.Data, req.Metadata)
	if err != nil {
		return &pb.ExecuteRuleResponse{
			Success: false,
			Result:  "",
			Message: err.Error(),
		}, nil
	}

	return &pb.ExecuteRuleResponse{
		Success: true,
		Result:  result,
		Message: "Rule executed successfully",
	}, nil
}

// PostMessage 推送消息到规则链
func (s *RuleService) PostMessage(ctx context.Context, req *pb.PostMessageRequest) (*pb.PostMessageResponse, error) {
	// 异步推送消息到规则链
	err := s.uc.PostMessage(ctx, req.Id, req.Username, req.MsgType, req.Data, req.Metadata)
	if err != nil {
		return &pb.PostMessageResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.PostMessageResponse{
		Success: true,
		Message: "Message posted successfully",
	}, nil
}

// OperateRule 操作规则链（部署/下线）
func (s *RuleService) OperateRule(ctx context.Context, req *pb.OperateRuleRequest) (*pb.OperateRuleResponse, error) {
	var err error
	switch req.Type {
	case "deploy":
		err = s.uc.DeployRule(ctx, req.Id, req.Username)
	case "undeploy":
		err = s.uc.UndeployRule(ctx, req.Id, req.Username)
	default:
		return &pb.OperateRuleResponse{
			Success: false,
			Message: "Invalid operation type: " + req.Type,
		}, nil
	}

	if err != nil {
		return &pb.OperateRuleResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.OperateRuleResponse{
		Success: true,
		Message: fmt.Sprintf("Rule %s completed successfully", req.Type),
	}, nil
}

// CreateComponentUseRule 创建组件使用规则
func (s *RuleService) CreateComponentUseRule(ctx context.Context, req *pb.CreateComponentUseRuleRequest) (*pb.CreateComponentUseRuleResponse, error) {
	// 简化实现，返回成功
	// 在实际项目中，这里应该保存到数据库
	s.log.WithContext(ctx).Infof("CreateComponentUseRule: %s for user %s", req.ComponentType, req.Username)

	return &pb.CreateComponentUseRuleResponse{
		Success: true,
		Message: "Component use rule created successfully",
		Id:      time.Now().Unix(), // 使用时间戳作为ID
	}, nil
}

// GetComponentUseRule 获取组件使用规则
func (s *RuleService) GetComponentUseRule(ctx context.Context, req *pb.GetComponentUseRuleRequest) (*pb.GetComponentUseRuleResponse, error) {
	// 简化实现，返回模拟数据
	// 在实际项目中，这里应该从数据库查询
	s.log.WithContext(ctx).Infof("GetComponentUseRule: %d for user %s", req.Id, req.Username)

	return &pb.GetComponentUseRuleResponse{
		Rule: &pb.ComponentUseRule{
			Id:            req.Id,
			ComponentType: "jsFilter",
			UsageInfo:     `{"description":"用于数据过滤的JavaScript组件","usage":"在规则中使用"}`,
			CreatedTime:   time.Now().Unix(),
			UpdatedTime:   time.Now().Unix(),
			Username:      req.Username,
		},
	}, nil
}

// DeleteComponentUseRule 删除组件使用规则
func (s *RuleService) DeleteComponentUseRule(ctx context.Context, req *pb.DeleteComponentUseRuleRequest) (*pb.DeleteComponentUseRuleResponse, error) {
	// 简化实现，返回成功
	// 在实际项目中，这里应该从数据库删除
	s.log.WithContext(ctx).Infof("DeleteComponentUseRule: %d for user %s", req.Id, req.Username)

	return &pb.DeleteComponentUseRuleResponse{
		Success: true,
		Message: "Component use rule deleted successfully",
	}, nil
}

// UpdateComponentUseRule 更新组件使用规则
func (s *RuleService) UpdateComponentUseRule(ctx context.Context, req *pb.UpdateComponentUseRuleRequest) (*pb.UpdateComponentUseRuleResponse, error) {
	// 简化实现，返回成功
	// 在实际项目中，这里应该更新数据库记录
	s.log.WithContext(ctx).Infof("UpdateComponentUseRule: %d for user %s", req.Id, req.Username)

	return &pb.UpdateComponentUseRuleResponse{
		Success: true,
		Message: "Component use rule updated successfully",
	}, nil
}

// ListComponentUseRules 分页查询组件使用规则
func (s *RuleService) ListComponentUseRules(ctx context.Context, req *pb.ListComponentUseRulesRequest) (*pb.ListComponentUseRulesResponse, error) {
	// 简化实现，返回模拟数据
	// 在实际项目中，这里应该从数据库分页查询
	s.log.WithContext(ctx).Infof("ListComponentUseRules for user %s", req.Username)

	// 模拟一些组件使用规则数据
	rules := []*pb.ComponentUseRule{
		{
			Id:            1,
			ComponentType: "jsFilter",
			UsageInfo:     `{"description":"用于数据过滤的JavaScript组件","example":"return msg.temperature > 25;"}`,
			CreatedTime:   time.Now().Unix(),
			UpdatedTime:   time.Now().Unix(),
			Username:      req.Username,
		},
		{
			Id:            2,
			ComponentType: "restApiCall",
			UsageInfo:     `{"description":"用于HTTP API调用的组件","example":"POST /api/data"}`,
			CreatedTime:   time.Now().Unix(),
			UpdatedTime:   time.Now().Unix(),
			Username:      req.Username,
		},
		{
			Id:            3,
			ComponentType: "log",
			UsageInfo:     `{"description":"用于输出日志的组件","level":"INFO,DEBUG,WARN,ERROR"}`,
			CreatedTime:   time.Now().Unix(),
			UpdatedTime:   time.Now().Unix(),
			Username:      req.Username,
		},
	}

	// 简单的类型过滤
	if req.ComponentType != "" {
		filteredRules := []*pb.ComponentUseRule{}
		for _, rule := range rules {
			if rule.ComponentType == req.ComponentType {
				filteredRules = append(filteredRules, rule)
			}
		}
		rules = filteredRules
	}

	return &pb.ListComponentUseRulesResponse{
		Rules:    rules,
		Total:    int32(len(rules)),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
