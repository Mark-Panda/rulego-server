package dao

import (
	"encoding/json"
	"sync"

	"github.com/rulego/rulego-server/config"
	"github.com/rulego/rulego/api/types"
)

type ComponentDao struct {
	config   config.Config
	username string
	index    Index
	sync.RWMutex
	ruleDao *RuleDao
}

func NewComponentDao(config config.Config, username string) (*ComponentDao, error) {

	ruleDao, err := NewRuleDaoToDataBase(config, username)
	if err != nil {
		return nil, err
	}
	dao := &ComponentDao{
		config:   config,
		username: username,
		index:    Index{Rules: make(map[string]RuleMeta)},
		ruleDao:  ruleDao,
	}
	return dao, nil
}
func (d *ComponentDao) List(username string, keywords string, root *bool, disabled *bool, size, page int) ([]types.RuleChain, int, error) {
	return d.ruleDao.ListToDataBase(username, keywords, root, disabled, size, page)
}

func (d *ComponentDao) Get(username, chainId string) ([]byte, error) {
	return d.ruleDao.FindDataBaseByRuleChainId(chainId)
}

func (d *ComponentDao) GetAsRuleChain(username, chainId string) (types.RuleChain, error) {
	// 根据ID加载规则链DSL数据
	var ruleChain types.RuleChain
	data, err := d.Get(username, chainId)
	if err != nil {
		return ruleChain, err
	}
	if err := json.Unmarshal(data, &ruleChain); err != nil {
		return ruleChain, err
	}

	return ruleChain, nil
}

func (d *ComponentDao) Save(username, chainId string, def []byte) error {
	var ruleChain types.RuleChain
	if err := json.Unmarshal(def, &ruleChain); err != nil {
		return err
	}
	if err := d.ruleDao.SaveToDataBase(username, chainId, def); err != nil {
		return err
	}
	return nil
}

// 从数据库删除规则链
func (d *ComponentDao) DeleteToDataBase(username, chainId string) error {
	return DeleteRegulationByRuleChainId(chainId)
}

func (d *ComponentDao) Delete(username, chainId string) error {
	return d.ruleDao.DeleteToDataBase(username, chainId)
}
