package controller

import (
	"net/http"
	"strconv"

	"github.com/rulego/rulego-server/internal/constants"
	"github.com/rulego/rulego-server/internal/model"
	"github.com/rulego/rulego-server/internal/service"
	endpointApi "github.com/rulego/rulego/api/types/endpoint"
	"github.com/rulego/rulego/endpoint"
	"github.com/rulego/rulego/utils/json"
)

var Doc = &doc{}

type doc struct {
}

// Create 创建业务md文档
func (d *doc) Create(url string) endpointApi.Router {
	return endpoint.NewRouter().From(url).Process(AuthProcess).Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		msg := exchange.In.GetMsg()
		var req model.MdWorkflow
		if err := json.Unmarshal([]byte(msg.GetData()), &req); err != nil {
			exchange.Out.SetStatusCode(http.StatusBadRequest)
			exchange.Out.SetBody([]byte(err.Error()))
		}
		if req.Title == "" || req.Content == "" || req.Desc == "" {
			exchange.Out.SetStatusCode(http.StatusBadRequest)
			exchange.Out.SetBody([]byte("title,content,desc不能为空"))
		} else {
			err := service.EventServiceImpl.CreateMdWorkflow(req)
			if err != nil {
				exchange.Out.SetStatusCode(http.StatusInternalServerError)
				exchange.Out.SetBody([]byte(err.Error()))
			}
		}
		return true
	}).End()
}

// Edit 编辑业务md文档
func (d *doc) Edit(url string) endpointApi.Router {
	return endpoint.NewRouter().From(url).Process(AuthProcess).Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		msg := exchange.In.GetMsg()
		var req model.MdWorkflow
		if err := json.Unmarshal([]byte(msg.GetData()), &req); err != nil {
			exchange.Out.SetStatusCode(http.StatusBadRequest)
			exchange.Out.SetBody([]byte(err.Error()))
		}
		if req.Id == 0 {
			exchange.Out.SetStatusCode(http.StatusBadRequest)
			exchange.Out.SetBody([]byte("数据未找到"))
		}
		updateInfo := map[string]interface{}{}
		if req.Title != "" {
			updateInfo["title"] = req.Title
		}
		if req.Desc != "" {
			updateInfo["desc"] = req.Desc
		}
		if req.Content != "" {
			updateInfo["content"] = req.Content
		}
		if req.ChainId != "" {
			updateInfo["chain_id"] = req.ChainId
		}
		if req.ChainName != "" {
			updateInfo["chain_name"] = req.ChainName
		}
		if req.ChainVersion != 0 {
			updateInfo["chain_version"] = req.ChainVersion
		}
		if err := service.EventServiceImpl.EditMdWorkflow(req.Id, updateInfo); err != nil {
			exchange.Out.SetStatusCode(http.StatusInternalServerError)
			exchange.Out.SetBody([]byte(err.Error()))
		}
		return true
	}).End()
}

// List 获取业务md文档列表
func (d *doc) List(url string) endpointApi.Router {
	return endpoint.NewRouter().From(url).Process(AuthProcess).Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		msg := exchange.In.GetMsg()
		title := msg.Metadata.GetValue("keyword")
		var current = 1
		var pageSize = 20
		var result interface{}
		currentStr := msg.Metadata.GetValue(constants.KeyPage)
		if i, err := strconv.Atoi(currentStr); err == nil {
			current = i
		}
		pageSizeStr := msg.Metadata.GetValue(constants.KeySize)
		if i, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = i
		}
		if v, total, err := service.EventServiceImpl.GetMdWorkflowList(title, current, pageSize); err != nil {
			exchange.Out.SetStatusCode(http.StatusNotFound)
			exchange.Out.SetBody([]byte(err.Error()))
			return false
		} else {
			result = map[string]interface{}{
				"page":  current,
				"size":  pageSize,
				"total": total,
				"items": v,
			}
		}
		if v, err := json.Marshal(result); err != nil {
			exchange.Out.SetStatusCode(http.StatusInternalServerError)
			exchange.Out.SetBody([]byte(err.Error()))
		} else {
			exchange.Out.SetBody(v)
		}
		return true
	}).End()
}
