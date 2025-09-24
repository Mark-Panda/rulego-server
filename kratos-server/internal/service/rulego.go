package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"kratos-server/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	endpointApi "github.com/rulego/rulego/api/types/endpoint"
	"github.com/rulego/rulego/endpoint"
	"github.com/rulego/rulego/endpoint/rest"
	"github.com/rulego/rulego/node_pool"

	// 导入RuleGo组件
	_ "github.com/rulego/rulego-components/endpoint/grpc_stream"
	_ "github.com/rulego/rulego-components/endpoint/kafka"
	_ "github.com/rulego/rulego-components/endpoint/nats"
	_ "github.com/rulego/rulego-components/endpoint/rabbitmq"
	_ "github.com/rulego/rulego-components/endpoint/redis"
	_ "github.com/rulego/rulego-components/stats/streamsql"
)

// RuleGoService RuleGo集成服务
type RuleGoService struct {
	config       *conf.RuleGo
	rulegoConfig types.Config
	nodePool     *node_pool.NodePool
	restEndpoint endpointApi.HttpEndpoint
	ruleEngines  map[string]types.RuleEngine // 存储规则引擎实例
	mu           sync.RWMutex                // 保护ruleEngines的并发访问
	log          *log.Helper
}

// NewRuleGoService 创建RuleGo集成服务
func NewRuleGoService(config *conf.Bootstrap, logger log.Logger) (*RuleGoService, error) {
	helper := log.NewHelper(logger)
	rulegoConf := config.GetRulego()

	// 初始化RuleGo配置
	rulegoConfig := rulego.NewConfig(types.WithDefaultPool())
	nodePool := node_pool.NewNodePool(rulegoConfig)
	rulegoConfig.NodePool = nodePool

	service := &RuleGoService{
		config:       rulegoConf,
		rulegoConfig: rulegoConfig,
		nodePool:     nodePool,
		ruleEngines:  make(map[string]types.RuleEngine),
		log:          helper,
	}

	// 初始化REST端点
	if err := service.initRestEndpoint(); err != nil {
		return nil, err
	}

	return service, nil
}

// initRestEndpoint 初始化REST端点
func (s *RuleGoService) initRestEndpoint() error {
	addr := ":9091" // 使用不同的端口避免与Kratos HTTP服务冲突

	s.log.Infof("RuleGo REST server starting at %s", addr)

	ep, err := endpoint.Registry.New(
		rest.Type,
		s.rulegoConfig,
		rest.Config{
			Server:    addr,
			AllowCors: true,
		},
	)
	if err != nil {
		return err
	}

	var restEndpoint endpointApi.HttpEndpoint
	if ep, ok := ep.(endpointApi.HttpEndpoint); !ok {
		return fmt.Errorf("is not HttpEndpoint type error")
	} else {
		restEndpoint = ep
	}

	// 添加全局拦截器
	restEndpoint.AddInterceptors(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		if out, ok := exchange.Out.(endpointApi.HeaderModifier); ok {
			out.AddHeader("Content-Type", "application/json")
		} else {
			exchange.Out.Headers().Set("Content-Type", "application/json")
		}
		return true
	})

	// 重定向UI界面
	restEndpoint.GET(endpoint.NewRouter().From("/").Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		r, ok1 := exchange.In.(*rest.RequestMessage)
		w, ok2 := exchange.Out.(*rest.ResponseMessage)
		if ok1 && ok2 {
			http.Redirect(w.Response(), r.Request(), "/editor/", http.StatusFound)
		}
		return false
	}).End())

	// 添加原RuleGo API路由
	// 这里可以集成原有的controller逻辑
	// 例如：规则链管理、组件管理、日志管理等API

	// 为了与Kratos HTTP服务区分，这里可以添加一些特定的RuleGo功能
	// 比如：直接通过RuleGo引擎执行的API端点

	// 示例：直接执行规则链的端点
	restEndpoint.POST(endpoint.NewRouter().From("/rulego/execute/:id").Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		// 这里可以直接通过RuleGo引擎执行规则链
		// 而不需要经过Kratos的gRPC/HTTP路由
		if out, ok := exchange.Out.(endpointApi.HeaderModifier); ok {
			out.AddHeader("Content-Type", "application/json")
		}
		// TODO: 实际的规则链执行逻辑
		exchange.Out.SetBody([]byte(`{"status":"success","message":"Rule executed via RuleGo directly"}`))
		return false
	}).End())

	// 加载静态文件映射
	if s.config.ResourceMapping != "" {
		restEndpoint.RegisterStaticFiles(s.config.ResourceMapping)
	}

	// 把HTTP服务设置成共享节点
	if s.config.ShareHttpServer {
		_, _ = node_pool.DefaultNodePool.AddNode(restEndpoint)
	}
	_, _ = s.nodePool.AddNode(restEndpoint)

	s.restEndpoint = restEndpoint
	return nil
}

// Start 启动RuleGo服务
func (s *RuleGoService) Start(ctx context.Context) error {
	if s.restEndpoint != nil {
		return s.restEndpoint.Start()
	}
	return nil
}

// Stop 停止RuleGo服务
func (s *RuleGoService) Stop(ctx context.Context) error {
	if s.restEndpoint != nil {
		s.restEndpoint.Destroy()
	}
	return nil
}

// GetNodePool 获取节点池
func (s *RuleGoService) GetNodePool() *node_pool.NodePool {
	return s.nodePool
}

// GetRestEndpoint 获取REST端点
func (s *RuleGoService) GetRestEndpoint() endpointApi.HttpEndpoint {
	return s.restEndpoint
}

// ExecuteRule 执行规则链
func (s *RuleGoService) ExecuteRule(ctx context.Context, ruleID, username, msgType, data string, metadata map[string]string) (string, error) {
	s.log.WithContext(ctx).Infof("ExecuteRule: ruleID=%s, username=%s, msgType=%s", ruleID, username, msgType)

	// 获取或创建规则引擎
	engine, err := s.getOrCreateRuleEngine(ruleID, username)
	if err != nil {
		return "", fmt.Errorf("failed to get rule engine: %w", err)
	}

	// 构造消息
	msg := types.NewMsg(0, msgType, types.JSON, types.NewMetadata(), data)

	// 添加元数据
	for k, v := range metadata {
		msg.Metadata.PutValue(k, v)
	}

	// 添加用户名到元数据
	msg.Metadata.PutValue("username", username)

	// 同步执行规则链并等待结果
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultChan := make(chan types.RuleMsg, 1)
	errorChan := make(chan error, 1)

	// 执行规则链
	engine.OnMsgAndWait(msg, types.WithContext(ctx), types.WithEndFunc(func(ctx types.RuleContext, msg types.RuleMsg, err error) {
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- msg
		}
	}))

	// 等待结果
	select {
	case result := <-resultChan:
		// 构造返回结果
		response := map[string]interface{}{
			"status":    "success",
			"data":      result.Data,
			"msgType":   result.Type,
			"timestamp": result.Ts,
			"metadata":  result.Metadata.Values(),
		}
		respBytes, _ := json.Marshal(response)
		return string(respBytes), nil
	case err := <-errorChan:
		return "", fmt.Errorf("rule execution failed: %w", err)
	case <-ctx.Done():
		return "", fmt.Errorf("rule execution timeout: %w", ctx.Err())
	}
}

// PostMessage 异步推送消息到规则链
func (s *RuleGoService) PostMessage(ctx context.Context, ruleID, username, msgType, data string, metadata map[string]string) error {
	s.log.WithContext(ctx).Infof("PostMessage: ruleID=%s, username=%s, msgType=%s", ruleID, username, msgType)

	// 获取或创建规则引擎
	engine, err := s.getOrCreateRuleEngine(ruleID, username)
	if err != nil {
		return fmt.Errorf("failed to get rule engine: %w", err)
	}

	// 构造消息
	msg := types.NewMsg(0, msgType, types.JSON, types.NewMetadata(), data)

	// 添加元数据
	for k, v := range metadata {
		msg.Metadata.PutValue(k, v)
	}

	// 添加用户名到元数据
	msg.Metadata.PutValue("username", username)

	// 异步执行规则链
	engine.OnMsg(msg, types.WithContext(ctx))

	return nil
}

// DeployRule 部署规则链
func (s *RuleGoService) DeployRule(ctx context.Context, ruleID, username, ruleDefinition string) error {
	s.log.WithContext(ctx).Infof("DeployRule: ruleID=%s, username=%s", ruleID, username)

	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果已存在，先停止旧的引擎
	if existingEngine, exists := s.ruleEngines[ruleID]; exists {
		existingEngine.Stop(ctx)
		delete(s.ruleEngines, ruleID)
	}

	// 创建新的规则引擎
	engine, err := s.createRuleEngine(ruleID, ruleDefinition)
	if err != nil {
		return fmt.Errorf("failed to create rule engine: %w", err)
	}

	// 存储引擎实例
	s.ruleEngines[ruleID] = engine

	s.log.WithContext(ctx).Infof("Rule deployed successfully: %s", ruleID)
	return nil
}

// UndeployRule 下线规则链
func (s *RuleGoService) UndeployRule(ctx context.Context, ruleID, username string) error {
	s.log.WithContext(ctx).Infof("UndeployRule: ruleID=%s, username=%s", ruleID, username)

	s.mu.Lock()
	defer s.mu.Unlock()

	// 查找并停止规则引擎
	if engine, exists := s.ruleEngines[ruleID]; exists {
		engine.Stop(ctx)
		delete(s.ruleEngines, ruleID)
		s.log.WithContext(ctx).Infof("Rule undeployed successfully: %s", ruleID)
		return nil
	}

	return fmt.Errorf("rule not found or not deployed: %s", ruleID)
}

// getOrCreateRuleEngine 获取或创建规则引擎
func (s *RuleGoService) getOrCreateRuleEngine(ruleID, username string) (types.RuleEngine, error) {
	s.mu.RLock()
	if engine, exists := s.ruleEngines[ruleID]; exists {
		s.mu.RUnlock()
		return engine, nil
	}
	s.mu.RUnlock()

	// 需要创建新的引擎，但这里需要规则定义
	// 在实际项目中，应该从数据库或其他存储中获取规则定义
	// 这里返回一个默认的示例规则链
	defaultRuleChain := `{
		"ruleChain": {
			"id": "` + ruleID + `",
			"name": "Default Rule Chain",
			"root": true
		},
		"metadata": {
			"nodes": [
				{
					"id": "s1",
					"type": "log",
					"name": "记录日志",
					"configuration": {
						"jsScript": "return '处理消息: ' + JSON.stringify(msg);"
					}
				}
			],
			"connections": [
				{
					"fromId": "root",
					"toId": "s1",
					"type": "Success"
				}
			]
		}
	}`

	return s.createRuleEngine(ruleID, defaultRuleChain)
}

// createRuleEngine 创建规则引擎
func (s *RuleGoService) createRuleEngine(ruleID, ruleDefinition string) (types.RuleEngine, error) {
	// 创建规则引擎实例
	engine, err := rulego.New(ruleID, []byte(ruleDefinition), rulego.WithConfig(s.rulegoConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to create rule engine for %s: %w", ruleID, err)
	}

	return engine, nil
}

// GetRuleEngine 获取指定的规则引擎（用于调试等场景）
func (s *RuleGoService) GetRuleEngine(ruleID string) (types.RuleEngine, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	engine, exists := s.ruleEngines[ruleID]
	return engine, exists
}

// ListDeployedRules 获取已部署的规则链列表
func (s *RuleGoService) ListDeployedRules() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ruleIDs := make([]string, 0, len(s.ruleEngines))
	for ruleID := range s.ruleEngines {
		ruleIDs = append(ruleIDs, ruleID)
	}
	return ruleIDs
}
