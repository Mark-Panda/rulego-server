package controller

import (
	"net/http"

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
		// service.EventServiceImpl.CreateMdWorkflow()
		return true
	}).End()
}

// Edit 编辑业务md文档
func (d *doc) Edit(url string) endpointApi.Router {
	return endpoint.NewRouter().From(url).Process(AuthProcess).Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {

		return true
	}).End()
}

// List 获取业务md文档列表
func (d *doc) List(url string) endpointApi.Router {
	return endpoint.NewRouter().From(url).Process(AuthProcess).Process(func(router endpointApi.Router, exchange *endpointApi.Exchange) bool {
		docInfo := map[string]string{
			"title":     "doc1",
			"desc":      "doc1 desc",
			"content":   "doc1 content",
			"chainId":   "1234567890",
			"createdAt": "2022-01-01 00:00:00",
			"updatedAt": "2022-01-01 00:00:00",
		}
		result := map[string]interface{}{
			"page":  1,
			"size":  10,
			"total": 1,
			"items": []map[string]string{docInfo},
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
