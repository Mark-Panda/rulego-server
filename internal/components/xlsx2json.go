package components

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/maps"
	"github.com/xuri/excelize/v2"
)

// init registers the component to rulego
func init() {
	_ = rulego.Registry.Register(&Xlsx2Json{})
}

// Xlsx2JsonConfiguration node configuration
type Xlsx2JsonConfiguration struct {
	// 文件下载链接
	DownLoadUrl string
}

type Xlsx2Json struct {
	Config Xlsx2JsonConfiguration
}

// New 创建一个组件新实例
// 每个规则链的规则节点都会创建一个新的实例，数据是隔离的
func (x *Xlsx2Json) New() types.Node {
	return &Xlsx2Json{Config: Xlsx2JsonConfiguration{
		DownLoadUrl: "",
	}}
}

// Type 组件类型，类型不能重复。
// 用于规则链，node.type配置，初始化对应的组件
// 建议使用`/`区分命名空间，防止冲突。例如：x/httpClient
func (x *Xlsx2Json) Type() string {
	return "transform/xlsx2json"
}

// 实现ComponentDefGetter接口修改组件名和描述
func (x *Xlsx2Json) Def() types.ComponentForm {
	relationTypes := &[]string{"Success", "Failure"}
	return types.ComponentForm{
		Label:         "xlsx2json",
		Desc:          "将xlsx文件转换为json",
		RelationTypes: relationTypes,
	}
}

// Init 组件初始化，一般做一些组件参数配置或者客户端初始化操作
// 规则链里的规则节点初始化会调用一次
func (x *Xlsx2Json) Init(ruleConfig types.Config, configuration types.Configuration) error {
	err := maps.Map2Struct(configuration, &x.Config)
	return err
}

// OnMsg 处理消息，并控制流向子节点的关系。每条流入组件的数据会经过该方法处理
// ctx:规则引擎处理消息上下文
// msg:消息
func (x *Xlsx2Json) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	// 逻辑处理
	if x.Config.DownLoadUrl == "" {
		ctx.TellFailure(msg, errors.New("文件地址未配置"))
	}
	// 下载并解析Excel文件
	data, err := downloadAndParseExcel(x.Config.DownLoadUrl)
	if err != nil {
		errMsg := fmt.Sprintf(`下载并解析Excel文件时出错:%s`, err.Error())
		ctx.TellFailure(msg, errors.New(errMsg))
	}

	// 将数据转换为JSON
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf(`转换为JSON时出错:%s`, err.Error())
		ctx.TellFailure(msg, errors.New(errMsg))
	}
	strJson := string(jsonBytes)
	msg.Data = strJson
	ctx.TellSuccess(msg)
}

// Destroy 销毁，做一些资源释放操作
func (x *Xlsx2Json) Destroy() {
	_ = x.Close()
}

func (x *Xlsx2Json) Close() error {
	return nil
}

// RowData 表示一行数据
type RowData map[string]interface{}

func downloadAndParseExcel(url string) ([]map[string][]RowData, error) {
	// 下载文件
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("无法下载文件，状态码: %d", resp.StatusCode)
	}
	// 读取文件内容到内存
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	// 从内存中的文件流读取Excel数据
	f, err := excelize.OpenReader(buf)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	allSheetMap := f.GetSheetMap()
	dataList := make([]map[string][]RowData, 0)
	for _, sheetName := range allSheetMap {
		dataMap := make(map[string][]RowData)
		// 读取指定的工作表
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return nil, err
		}
		var header []string
		var data []RowData
		// 处理标题行
		header = append(header, rows[0]...)
		// 处理数据行
		for _, row := range rows[1:] {
			rowData := make(RowData)
			for i, cell := range row {
				if i < len(header) {
					if isNumber(cell) {
						if intValue, err := strconv.Atoi(cell); err == nil {
							rowData[header[i]] = intValue
						} else if floatValue, err := strconv.ParseFloat(cell, 64); err == nil {
							rowData[header[i]] = floatValue
						} else {
							rowData[header[i]] = cell
						}
					} else {
						rowData[header[i]] = cell
					}
				}
			}
			data = append(data, rowData)
		}
		dataMap[sheetName] = data
		dataList = append(dataList, dataMap)
	}
	return dataList, nil
}

// isNumber 检查字符串是否为数字
func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
