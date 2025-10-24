package service

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/rulego/rulego-server/internal/model"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/json"
)

// 创建Md工作流
func (s *EventService) CreateMdWorkflow(info model.MdWorkflow) error {
	return s.EventDao.SaveMdWorkFlow(info)
}

// 编辑Md工作流
func (s *EventService) EditMdWorkflow(mdId int64, updateInfo map[string]interface{}) error {
	return s.EventDao.EditorMdWorkFlow(mdId, updateInfo)
}

// 分页获取md工作流列表
func (s *EventService) GetMdWorkflowList(title string, current, size int) ([]*model.MdWorkflow, int64, error) {
	return s.EventDao.GetMdWorkflowList(title, current, size)
}

func (s *EventService) GenerateWorkflow(mdId string, chainId string) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var result string
	ruleMsg := types.NewMsgWithJsonData("{\"msgType\":\"system_generate\"}")
	err := s.ruleEngineService.ExecuteAndWait(mdId, ruleMsg, types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
		result = msg.GetData()
		wg.Done()
	}))
	if err != nil {
		return err
	}
	wg.Wait()
	// DSL转化
	var ruleChain types.RuleChain
	err = json.Unmarshal([]byte(result), &ruleChain)
	if err != nil {
		return errors.New("规则链解析异常")
	}
	if chainId == "" {
		// 生成新的UUID
		chainId = uuid.New().String()
	}
	ruleChain.RuleChain.ID = chainId
	// 循环替换生成新的节点ID
	for i := range ruleChain.Metadata.Nodes {
		oldId := ruleChain.Metadata.Nodes[i].Id
		newId := uuid.New().String()
		ruleChain.Metadata.Nodes[i].Id = newId
		for j := range ruleChain.Metadata.Connections {
			if ruleChain.Metadata.Connections[j].FromId == oldId {
				ruleChain.Metadata.Connections[j].FromId = newId
			}
			if ruleChain.Metadata.Connections[j].ToId == oldId {
				ruleChain.Metadata.Connections[j].ToId = newId
			}
		}
	}
	// 保存规则链
	b, err := json.Marshal(ruleChain)
	if err != nil {
		return err
	}
	if err = s.ruleEngineService.SaveAndLoad(chainId, b); err != nil {
		return err
	}
	return nil
}
