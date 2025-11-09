package components

import (
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/maps"
)

// init registers the component to rulego
func init() {
	_ = rulego.Registry.Register(&StartDsl{})
}

type StartDsl struct {
	// 节点配置
	Config StartDslConfiguration
}

type StartDslConfiguration struct {
}

func (c *StartDsl) New() types.Node {
	return &StartDsl{Config: StartDslConfiguration{}}
}

// Type 组件类型，类型不能重复。
// 用于规则链，node.type配置，初始化对应的组件
// 建议使用`/`区分命名空间，防止冲突。例如：x/httpClient
func (c *StartDsl) Type() string {
	return "start"
}

// 实现ComponentDefGetter接口修改组件名和描述
func (c *StartDsl) Def() types.ComponentForm {
	// relationTypes := &[]string{"Success", "Failure"}
	return types.ComponentForm{
		Label: "start",
		Desc:  "虚拟开始节点",
		// RelationTypes: relationTypes,
	}
}

// Init 组件初始化，一般做一些组件参数配置或者客户端初始化操作
// 规则链里的规则节点初始化会调用一次
func (c *StartDsl) Init(ruleConfig types.Config, configuration types.Configuration) error {
	err := maps.Map2Struct(configuration, &c.Config)
	return err
}

// OnMsg 处理消息，并控制流向子节点的关系。每条流入组件的数据会经过该方法处理
// ctx:规则引擎处理消息上下文
// msg:消息
func (c *StartDsl) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {

	ctx.TellSuccess(msg)
}

// Destroy 销毁，做一些资源释放操作
func (c *StartDsl) Destroy() {
	_ = c.Close()
}

func (c *StartDsl) Close() error {
	return nil
}
