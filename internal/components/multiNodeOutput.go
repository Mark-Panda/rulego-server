package components

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/components/base"
	"github.com/rulego/rulego/utils/maps"
)

// init registers the component to rulego
func init() {
	_ = rulego.Registry.Register(&MultiNodeOutput{})
}

type MultiNodeOutput struct {
	// 节点配置
	Config MultiNodeOutputConfiguration
}

type MultiNodeOutputConfiguration struct {
	// NodeId 目标节点ID，获取该节点的输出消息
	NodeIds []string `json:"nodeIds"`
}

func (c *MultiNodeOutput) New() types.Node {
	return &MultiNodeOutput{Config: MultiNodeOutputConfiguration{}}
}

// Type 组件类型，类型不能重复。
// 用于规则链，node.type配置，初始化对应的组件
// 建议使用`/`区分命名空间，防止冲突。例如：x/httpClient
func (c *MultiNodeOutput) Type() string {
	return "transform/multiNodeOutput"
}

// 实现ComponentDefGetter接口修改组件名和描述
func (c *MultiNodeOutput) Def() types.ComponentForm {
	return types.ComponentForm{
		Label: "multiNodeOutput",
		Desc:  "获取已完成节点的输出信息",
	}
}

// Init 组件初始化，一般做一些组件参数配置或者客户端初始化操作
// 规则链里的规则节点初始化会调用一次
func (c *MultiNodeOutput) Init(ruleConfig types.Config, configuration types.Configuration) error {
	err := maps.Map2Struct(configuration, &c.Config)
	if err != nil {
		return err
	}
	chainCtx := base.NodeUtils.GetChainCtx(configuration)
	if chainCtx == nil {
		return errors.New("chain ctx is nil")
	}
	self := base.NodeUtils.GetSelfDefinition(configuration)
	// Establish node dependency to enable target node output caching and access
	for _, nodeId := range c.Config.NodeIds {
		chainCtx.AddNodeDependency(self.Id, nodeId)
	}
	return err
}

// OnMsg 处理消息，并控制流向子节点的关系。每条流入组件的数据会经过该方法处理
// ctx:规则引擎处理消息上下文
// msg:消息
func (c *MultiNodeOutput) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	resultData := map[string]interface{}{}
	isSuccess := true
	errStr := ""
	for _, nodeId := range c.Config.NodeIds {
		if targetMsg, exists := ctx.GetNodeRuleMsg(nodeId); exists {
			// 合并多节点返回数据
			resultData[nodeId] = targetMsg.GetData()
		} else {
			isSuccess = false
			errStr += fmt.Sprintf("node %s output not found; ", nodeId)
		}
	}
	if isSuccess {
		// 合并获取后，发送到成功链
		resultBytes, _ := json.Marshal(resultData)
		msg.SetData(string(resultBytes))
		ctx.TellSuccess(msg)
	} else {
		// Target node has no output or dependency not established, send to failure chain
		ctx.TellFailure(msg, errors.New(errStr))
	}
}

// Destroy 销毁，做一些资源释放操作
func (c *MultiNodeOutput) Destroy() {
	_ = c.Close()
}

func (c *MultiNodeOutput) Close() error {
	return nil
}
