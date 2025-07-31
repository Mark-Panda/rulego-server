package service

import (
	"github.com/rulego/rulego-server/config"
	"github.com/rulego/rulego-server/internal/dao"
	"github.com/rulego/rulego-server/internal/model"
	"github.com/rulego/rulego/api/types"
)

var EventServiceImpl *EventService

type EventService struct {
	EventDao *dao.EventDao
	config   config.Config
}

func NewEventService(config config.Config) (*EventService, error) {
	if eventDao, err := dao.NewEventDao(config); err != nil {
		return nil, err
	} else {
		return &EventService{
			EventDao: eventDao,
			config:   config,
		}, nil
	}
}

// SaveRunLog 保存工作流运行日志快照
func (s *EventService) SaveRunLog(username string, ctx types.RuleContext, snapshot types.RuleChainRunSnapshot) error {
	//return s.EventDao.SaveRunLog(username, ctx, snapshot)
	return s.EventDao.SaveRunLogToDataBase(username, ctx, snapshot)
}

func (s *EventService) Delete(username, chainId, id string) error {
	// return s.EventDao.Delete(username, chainId, id)
	return s.EventDao.DeleteDataBaseByRunId(username, chainId, id)
}
func (s *EventService) DeleteByChainId(username, chainId string) error {
	// return s.EventDao.DeleteByChainId(username, chainId)
	return s.EventDao.DeleteDataBaseByChainId(username, chainId)
}

func (s *EventService) List(username, chainId string, current, size int, startTime, endTime string) ([]types.RuleChainRunSnapshot, int, error) {
	//return s.EventDao.List(username, chainId, current, size)
	return s.EventDao.ListByDataBase(username, chainId, current, size, startTime, endTime)
}

func (s *EventService) Get(username, chainId, snapshotId string) (types.RuleChainRunSnapshot, error) {
	//return s.EventDao.Get(username, chainId, snapshotId)
	return s.EventDao.GetByDataBase(username, chainId, snapshotId)
}

// CreateComponentUseRule 创建组件使用规则
func (s *EventService) CreateComponentUseRule(r model.ComponentUseRule) error {
	return s.EventDao.CreateComponentUseRule(r)
}

// FindComponentUseRuleByPage 分页查询组件使用规则
func (s *EventService) FindComponentUseRuleByPage(page, size int, component_type string, disabled *bool, keywords string) ([]model.ComponentUseRule, int64, error) {
	return s.EventDao.FindComponentUseRuleByPage(page, size, component_type, disabled, keywords)
}

// DeleteComponentUseRuleById 删除组件使用规则
func (s *EventService) DeleteComponentUseRuleById(id string) error {
	return s.EventDao.DeleteComponentUseRuleByIdPhysical(id)
}

// UpdateComponentUseRuleById 根据ID更新组件使用规则
func (s *EventService) UpdateComponentUseRuleById(id string, data map[string]interface{}) error {
	return s.EventDao.UpdateComponentUseRuleById(id, data)
}

// FindComponentUseRuleById 根据ID查询组件使用规则
func (s *EventService) FindComponentUseRuleById(id string) (model.ComponentUseRule, error) {
	return s.EventDao.FindComponentUseRuleById(id)
}
