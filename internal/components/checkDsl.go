package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/components/base"
	"github.com/rulego/rulego/utils/maps"
	"github.com/rulego/rulego/utils/str"
)

// init registers the component to rulego
func init() {
	_ = rulego.Registry.Register(&CheckDsl{})
}

type CheckDsl struct {
	// 节点配置
	Config CheckDslConfiguration
	// 规则是否包含变量
	ruleHasVar bool
}

type CheckDslConfiguration struct {
	// 规则链
	// 支持${metadata.key}占位符读取metadata元数据
	// 支持${msg.key}占位符读取消息负荷指定key数据
	// 支持${data}获取消息原始负荷
	Rule string
	// 是否启用详细验证
	EnableDetailValidation bool
	// 是否验证节点类型
	ValidateNodeTypes bool
}

func (c *CheckDsl) New() types.Node {
	return &CheckDsl{Config: CheckDslConfiguration{
		Rule:                   "",
		EnableDetailValidation: true,
		ValidateNodeTypes:      true,
	}}
}

// Type 组件类型，类型不能重复。
// 用于规则链，node.type配置，初始化对应的组件
// 建议使用`/`区分命名空间，防止冲突。例如：x/httpClient
func (c *CheckDsl) Type() string {
	return "transform/checkdsl"
}

// 实现ComponentDefGetter接口修改组件名和描述
func (c *CheckDsl) Def() types.ComponentForm {
	// relationTypes := &[]string{"Success", "Failure"}
	return types.ComponentForm{
		Label: "checkdsl",
		Desc:  "检查是否符合dsl规则",
		// RelationTypes: relationTypes,
	}
}

// Init 组件初始化，一般做一些组件参数配置或者客户端初始化操作
// 规则链里的规则节点初始化会调用一次
func (c *CheckDsl) Init(ruleConfig types.Config, configuration types.Configuration) error {
	err := maps.Map2Struct(configuration, &c.Config)
	if err == nil {
		if str.CheckHasVar(c.Config.Rule) {
			c.ruleHasVar = true
		}
	}
	return err
}

// OnMsg 处理消息，并控制流向子节点的关系。每条流入组件的数据会经过该方法处理
// ctx:规则引擎处理消息上下文
// msg:消息
func (c *CheckDsl) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	var evn map[string]interface{}
	if c.ruleHasVar {
		evn = base.NodeUtils.GetEvnAndMetadata(ctx, msg)
	}
	ruleStr := c.Config.Rule
	if c.ruleHasVar {
		ruleStr = str.ExecuteTemplate(c.Config.Rule, evn)
	}
	if ruleStr == "" {
		ctx.TellFailure(msg, errors.New("未检查到规则链输入"))
		return
	}

	// 执行DSL格式检查
	validationResult, err := c.validateDsl(ruleStr)
	if err != nil {
		ctx.TellFailure(msg, err)
		return
	}

	// 将验证结果和原始DSL一起传递
	resultData := map[string]interface{}{
		"dsl":        ruleStr,
		"validation": validationResult,
		"isValid":    validationResult.IsValid,
	}

	resultBytes, _ := json.Marshal(resultData)
	msg.SetData(string(resultBytes))

	if validationResult.IsValid {
		ctx.TellSuccess(msg)
	} else {
		ctx.TellFailure(msg, fmt.Errorf("DSL验证失败: %s", strings.Join(validationResult.Errors, "; ")))
	}
}

// DSL验证结果结构
type DslValidationResult struct {
	IsValid   bool     `json:"isValid"`
	Errors    []string `json:"errors"`
	Warnings  []string `json:"warnings"`
	NodeCount int      `json:"nodeCount"`
}

// validateDsl 验证DSL格式
func (c *CheckDsl) validateDsl(ruleStr string) (*DslValidationResult, error) {
	result := &DslValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// 基础JSON格式检查
	var rule types.RuleChain
	if err := json.Unmarshal([]byte(ruleStr), &rule); err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("JSON格式错误: %v", err))
		return result, nil
	}

	// 如果启用详细验证
	if c.Config.EnableDetailValidation {
		c.performDetailedValidation(&rule, result)
	}

	return result, nil
}

// performDetailedValidation 执行详细的DSL验证
func (c *CheckDsl) performDetailedValidation(rule *types.RuleChain, result *DslValidationResult) {
	// 验证规则链基本信息
	if rule.RuleChain.ID == "" {
		result.Errors = append(result.Errors, "规则链ID不能为空")
		result.IsValid = false
	}

	if rule.RuleChain.Name == "" {
		result.Warnings = append(result.Warnings, "规则链名称为空")
	}

	// 验证节点信息
	nodeIds := make(map[string]bool)
	result.NodeCount = len(rule.Metadata.Nodes)

	for i, node := range rule.Metadata.Nodes {
		// 检查节点ID
		if node.Id == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("节点[%d]的ID不能为空", i))
			result.IsValid = false
			continue
		}

		// 检查ID重复
		if nodeIds[node.Id] {
			result.Errors = append(result.Errors, fmt.Sprintf("节点ID重复: %s", node.Id))
			result.IsValid = false
		}
		nodeIds[node.Id] = true

		// 检查节点类型
		if node.Type == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("节点[%s]的类型不能为空", node.Id))
			result.IsValid = false
		}

		// 验证节点类型是否合法（如果启用）
		if c.Config.ValidateNodeTypes {
			c.validateNodeType(node, result)
		}

		// 检查节点名称
		if node.Name == "" {
			result.Warnings = append(result.Warnings, fmt.Sprintf("节点[%s]的名称为空", node.Id))
		}
	}

	// 验证连接关系
	c.validateConnections(rule, nodeIds, result)
}

// validateNodeType 验证节点类型是否合法
func (c *CheckDsl) validateNodeType(node *types.RuleNode, result *DslValidationResult) {
	// 常见的节点类型列表（可以根据实际需要扩展）
	validTypes := map[string]bool{
		// 基础组件
		"filter":    true,
		"transform": true,
		"action":    true,
		"flow":      true,
		"external":  true,
		"start":     true,
		// 扩展组件
		"jsFilter":      true,
		"jsTransform":   true,
		"restApiCall":   true,
		"mqttClient":    true,
		"dbClient":      true,
		"log":           true,
		"delay":         true,
		"groupAction":   true,
		"groupFilter":   true,
		"switch":        true,
		"msgTypeSwitch": true,
		// 自定义组件
		"transform/checkdsl": true,
	}

	// 检查是否包含命名空间
	if strings.Contains(node.Type, "/") {
		parts := strings.Split(node.Type, "/")
		if len(parts) >= 2 {
			// 有命名空间的组件，暂时认为合法
			return
		}
	}

	if !validTypes[node.Type] {
		result.Warnings = append(result.Warnings, fmt.Sprintf("未知的节点类型: %s (节点: %s)", node.Type, node.Id))
	}
}

// validateConnections 验证节点连接关系
func (c *CheckDsl) validateConnections(rule *types.RuleChain, nodeIds map[string]bool, result *DslValidationResult) {
	// 查找起始节点
	startNodes := 0
	for _, node := range rule.Metadata.Nodes {
		if node.Type == "start" {
			startNodes++
		}
	}

	if startNodes == 0 {
		result.Warnings = append(result.Warnings, "未找到起始节点")
	} else if startNodes > 1 {
		result.Warnings = append(result.Warnings, "存在多个起始节点")
	}

	// 验证连接中引用的节点是否存在
	for _, connection := range rule.Metadata.Connections {
		if !nodeIds[connection.FromId] {
			result.Errors = append(result.Errors, fmt.Sprintf("连接中引用了不存在的源节点: %s", connection.FromId))
			result.IsValid = false
		}
		if !nodeIds[connection.ToId] {
			result.Errors = append(result.Errors, fmt.Sprintf("连接中引用了不存在的目标节点: %s", connection.ToId))
			result.IsValid = false
		}
	}
}

// Destroy 销毁，做一些资源释放操作
func (c *CheckDsl) Destroy() {
	_ = c.Close()
}

func (c *CheckDsl) Close() error {
	return nil
}
