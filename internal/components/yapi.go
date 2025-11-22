package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/maps"
)

// init registers the component to rulego
func init() {
	_ = rulego.Registry.Register(&YapiDsl{})
}

type YapiDsl struct {
	// 节点配置
	Config YapiDslConfiguration
}

type YapiDslConfiguration struct {
	BaseURL       string `json:"baseUrl"`
	LoginType     string `json:"loginType"`
	UserName      string `json:"userName"`
	Password      string `json:"password"`
	InterfacePath string `json:"interfacePath"`
}

func (c *YapiDsl) New() types.Node {
	return &YapiDsl{Config: YapiDslConfiguration{
		BaseURL:       "https://yapi.xxx.tv",
		LoginType:     "ldap",
		UserName:      "",
		Password:      "",
		InterfacePath: "",
	}}
}

// Type 组件类型，类型不能重复。
// 用于规则链，node.type配置，初始化对应的组件
// 建议使用`/`区分命名空间，防止冲突。例如：x/httpClient
func (c *YapiDsl) Type() string {
	return "transform/yapi"
}

// 实现ComponentDefGetter接口修改组件名和描述
func (c *YapiDsl) Def() types.ComponentForm {
	// relationTypes := &[]string{"Success", "Failure"}
	return types.ComponentForm{
		Label: "yapi",
		Desc:  "yapi接口获取",
		// RelationTypes: relationTypes,
	}
}

// Init 组件初始化，一般做一些组件参数配置或者客户端初始化操作
// 规则链里的规则节点初始化会调用一次
func (c *YapiDsl) Init(ruleConfig types.Config, configuration types.Configuration) error {
	err := maps.Map2Struct(configuration, &c.Config)
	return err
}

// OnMsg 处理消息，并控制流向子节点的关系。每条流入组件的数据会经过该方法处理
// ctx:规则引擎处理消息上下文
// msg:消息
func (c *YapiDsl) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {

	// 1. 从链接中提取接口 ID
	interfaceID, err := parseInterfaceID(c.Config.InterfacePath)
	if err != nil {
		fmt.Println("错误: 提取接口ID失败:", err)
		ctx.TellFailure(msg, err)
		return
	}
	client := resty.New().R()
	// 2. 执行登录操作
	cookies, err := performLogin(client, c.Config)
	if err != nil {
		fmt.Println("错误: 登录失败:", err)
		ctx.TellFailure(msg, err)
		return
	}
	// 3. 获取接口信息
	interfaceInfo, err := getInterfaceInfo(client, c.Config.BaseURL, interfaceID, cookies)
	if err != nil {
		fmt.Println("错误: 获取接口信息失败:", err)
		ctx.TellFailure(msg, err)
		return
	}
	jsonBytes, err := json.Marshal(interfaceInfo)
	if err != nil {
		fmt.Println("错误: JSON序列化失败:", err)
		ctx.TellFailure(msg, err)
		return
	}
	fmt.Println("接口信息:", interfaceInfo)
	msg.SetData(string(jsonBytes))
	ctx.TellSuccess(msg)
}

// Destroy 销毁，做一些资源释放操作
func (c *YapiDsl) Destroy() {
	_ = c.Close()
}

func (c *YapiDsl) Close() error {
	return nil
}

// parseInterfaceID 从 YAPI 接口链接中提取接口 ID
// 示例链接: http://yapi.example.com/project/123/interface/api/456
func parseInterfaceID(interfaceLink string) (int, error) {
	// 假设接口链接的末尾或者某个特定参数是接口 ID
	// 实际 YAPI 链接中，ID 通常在 URL 路径或查询参数中

	// 简单的处理方法是直接从链接中提取 ID（可能不健壮）
	u, err := url.Parse(interfaceLink)
	if err != nil {
		return 0, fmt.Errorf("解析接口链接 URL 失败: %w", err)
	}

	// 1) 先尝试查询参数，如 /api/interface/get?id=130644 或 ?interface_id=130644
	q := u.Query()
	if idStr := q.Get("id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil && id > 0 {
			return id, nil
		}
	}
	if idStr := q.Get("interface_id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil && id > 0 {
			return id, nil
		}
	}

	// 2) 再尝试从路径中提取
	// 常见形式：/project/{pid}/interface/api/{id}
	// 有时也可能是 /interface/api/{id}
	path := strings.Trim(u.Path, "/")
	parts := strings.Split(path, "/")

	// 2.1) 找到 "interface" 后紧跟 "api" 的下一个段为 id
	for i := 0; i < len(parts); i++ {
		if parts[i] == "interface" && i+2 < len(parts) && parts[i+1] == "api" {
			if id, err := strconv.Atoi(parts[i+2]); err == nil && id > 0 {
				return id, nil
			}
		}
	}

	// 2.2) 兜底：若最后一个段是纯数字，则认为是 id（例如路径末尾直接是 id）
	if len(parts) > 0 {
		last := parts[len(parts)-1]
		if id, err := strconv.Atoi(last); err == nil && id > 0 {
			return id, nil
		}
	}

	return 0, fmt.Errorf("未能从链接中解析到有效的接口ID: %s", interfaceLink)
}

// performLogin 执行登录操作，返回 http.Cookie 列表（包含登录信息）
func performLogin(client *resty.Request, config YapiDslConfiguration) ([]*http.Cookie, error) {
	if config.LoginType == "" {
		return nil, fmt.Errorf("登录类型不能为空")
	}
	loginPath := "/api/user/login"
	if config.LoginType == "ldap" {
		loginPath = "/api/user/login_by_ldap"
	}
	loginURL := config.BaseURL + loginPath

	resp, err := client.
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"email":    config.UserName,
			"password": config.Password,
		}).
		Post(loginURL)
	if err != nil {
		return nil, fmt.Errorf("登录请求失败: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("登录失败，状态码: %d, 响应: %s", resp.StatusCode(), string(resp.Body()))
	}

	var loginRes struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(bytes.NewReader(resp.Body())).Decode(&loginRes); err != nil {
		return nil, fmt.Errorf("解析登录响应失败: %w", err)
	}
	if loginRes.ErrCode != 0 {
		return nil, fmt.Errorf("YAPI 登录返回错误: %s", loginRes.ErrMsg)
	}

	return resp.Cookies(), nil
}

func getInterfaceInfo(client *resty.Request, baseURL string, interfaceID int, cookies []*http.Cookie) (interface{}, error) {
	interfacePath := fmt.Sprintf("/api/interface/get?id=%d", interfaceID)
	interfaceURL := baseURL + interfacePath

	resp, err := client.
		SetCookies(cookies).
		Get(interfaceURL)
	if err != nil {
		return nil, fmt.Errorf("获取接口信息请求失败: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("获取接口信息失败，状态码: %d, 响应: %s", resp.StatusCode(), string(resp.Body()))
	}
	var yapiRes YapiInterfaceResponse
	bodyBytes := resp.Body()
	fmt.Println("------", string(bodyBytes))
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&yapiRes); err != nil {
		return nil, fmt.Errorf("解析接口信息 JSON 失败: %w. 原始响应: %s", err, string(bodyBytes))
	}
	if yapiRes.ErrCode != 0 || yapiRes.Data == nil {
		return nil, fmt.Errorf("YAPI 接口返回错误: %s", yapiRes.ErrMsg)
	}

	return yapiRes.Data, nil
}

// 接口信息的 JSON 结构体，根据 YAPI 实际返回的结构定义，这里只定义了关键字段
type YapiInterfaceResponse struct {
	ErrCode int         `json:"errcode"`
	ErrMsg  string      `json:"errmsg"`
	Data    interface{} `json:"data"` // Data字段包含具体的接口信息
}
