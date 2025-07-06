package components

import (
	"encoding/json"
	"errors"

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
}

func (c *CheckDsl) New() types.Node {
	return &CheckDsl{Config: CheckDslConfiguration{
		Rule: "",
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
	var rule types.RuleChain
	err := json.Unmarshal([]byte(ruleStr), &rule)
	if err != nil {
		ctx.TellFailure(msg, errors.New("dsl规则解析失败"))
		return
	}
	msg.SetData(ruleStr)
	ctx.TellSuccess(msg)
}

// Destroy 销毁，做一些资源释放操作
func (c *CheckDsl) Destroy() {
	_ = c.Close()
}

func (c *CheckDsl) Close() error {
	return nil
}
