package constants

const (
	// DirWorkflows 工作流目录
	DirWorkflows          = "workflows"
	DirLocales            = "locales"
	DirWorkflowsRun       = "runs"
	DirWorkflowsRule      = "rules"
	DirWorkflowsComponent = "components"
	// DirWorkflowsShareNodes 共享节点目录
	DirWorkflowsShareNodes = "sharenodes"
	// DirShareNode 共享普通节点目录
	DirShareNode = "node"
	// DirShareEndpoint 共享端点目录
	DirShareEndpoint = "endpoint"
	// FileNameIndex 索引文件名
	FileNameIndex = "index"
)

const (
	// TypeShareNode 共享普通节点类型
	TypeShareNode = DirShareNode
	// TypeShareEndpoint 共享端点类型
	TypeShareEndpoint = DirShareEndpoint
)

// GetShareNodesDir 获取共享节点目录
func GetShareNodesDir() []string {
	return []string{
		DirShareNode,
		DirShareEndpoint,
	}
}

const (
	KeyMsgType         = "msgType"
	KeyMsgId           = "msgId"
	KeyChainId         = "chainId"
	KeyNodeId          = "nodeId"
	KeyUsername        = "username"
	KeyClientId        = "clientId"
	KeyVarType         = "varType"
	KeySize            = "size"
	KeyPage            = "page"
	KeyId              = "id"
	KeyKeywords        = "keywords"
	KeyType            = "type"
	KeyLang            = "lang"
	KeyRoot            = "root"
	KeyDisabled        = "disabled"
	KeyWebhookSecret   = "webhookSecret"
	KeyIntegrationType = "integrationType"
	// KeyWorkDir 工作目录
	KeyWorkDir = "workDir"
	// KeyDefaultIntegrationChainId 应用集成规则链ID
	KeyDefaultIntegrationChainId = "$event_bus"
	KeyUpdateTime                = "updateTime"
	KeyHeadersToMetadata         = "headersToMetadata"
	KeyInMessage                 = "inMessage"
	KeyBody                      = "body"
	KeyStartTime                 = "startTime"
	KeyEndTime                   = "endTime"
)

const (
	// OperateDeploy 部署
	OperateDeploy = "start"
	// OperateUndeploy 下架
	OperateUndeploy = "stop"
	// OperateSetToMain 设置成主规则链
	OperateSetToMain = "set-to-main"
)
const (
	// SettingKeyLatestChainId 最新打开的规则链
	SettingKeyLatestChainId = "latestChainId"
	// SettingKeyMainChainId 主规则链，server所有事件都会发送至此
	SettingKeyMainChainId = "mainChainId"
)

const (
	UserSuper = "super"
	UserAdmin = "admin"
)
const (
	RuleChainFileSuffix = ".json"
)
const (
	// AddiKeyMessage 记录规则链加载错误，扩展字段错误信息Key
	AddiKeyMessage = "message"
)
const (
	KeyAuthorization = "Authorization"
	KeyBearer        = "Bearer "
)

// LoadLuaLibs 加载lua库key
const LoadLuaLibs = "load_lua_libs"

// const (
//	DefaultPoolDef = `
//	{
//	  "ruleChain": {
//		"id": "$default_node_pool",
//		"name": "全局共享节点池"
//	  },
//	  "metadata": {
//		"endpoints": [
//		  {
//			"id": "core_endpoint_http",
//			"type": "endpoint/http",
//			"name": "http:9090",
//			"configuration": {
//			  "allowCors": true,
//			  "server": ":9090"
//			}
//		  }
//		]
//	  }
//	}
// `
// )
