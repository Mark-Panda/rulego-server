package dao

import (
	"sync"

	"github.com/dromara/carbon/v2"
	"github.com/rulego/rulego-server/config"
	"github.com/rulego/rulego-server/internal/model"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/json"
)

// IndexKeySpe key 连接符
var IndexKeySpe = ":"

type RuleDao struct {
	config   config.Config
	username string
	index    Index
	sync.RWMutex
}

// Index 定义索引结构，仅包含必要元数据
type Index struct {
	// key=chainId
	Rules map[string]RuleMeta `json:"rules"`
}

type RuleMeta struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Root       bool   `json:"root"`
	Disabled   bool   `json:"disabled"`
	UpdateTime string `json:"updateTime"`
}

// 保存或更新到数据库
func (d *RuleDao) SaveToComponentRegulation(username, chainId string, def []byte) error {
	v, _ := json.Format(def)
	// def 转成 types.RuleChain
	var ruleChain types.RuleChain
	if err := json.Unmarshal(def, &ruleChain); err != nil {
		return err
	}
	t := carbon.Now(carbon.Shanghai).StdTime()
	ruleConfigInfo, gErr := FindComponentRegulationByRuleChainId(chainId)
	if gErr != nil {
		return gErr
	}
	if ruleConfigInfo != nil && ruleConfigInfo.RuleChainId != "" {
		updateData := map[string]interface{}{
			"rule_config": string(v),
			"root":        ruleChain.RuleChain.Root,
			"disabled":    ruleChain.RuleChain.Disabled,
			"name":        ruleChain.RuleChain.Name,
			"updated_at":  &t,
		}

		return UpdateComponentRegulationByRuleChainId(chainId, updateData)
	}

	createInfo := model.ComponentRegulation{
		UserName:    username,
		Root:        ruleChain.RuleChain.Root,
		Disabled:    ruleChain.RuleChain.Disabled,
		Name:        ruleChain.RuleChain.Name,
		RuleChainId: chainId,
		RuleConfig:  string(v),
		CreatedAt:   &t,
		UpdatedAt:   &t,
	}
	return CreateComponentRegulation(createInfo)
}

// 从数据库删除规则链
func (d *RuleDao) DeleteToComponentRegulation(username, chainId string) error {
	return DeleteComponentRegulationByRuleChainIdPhysical(chainId)
}

// 按规则链id从数据库查询规则链
func (d *RuleDao) FindComponentRegulationByRuleChainId(chainId string) ([]byte, error) {
	ruleConfigInfo, gErr := FindComponentRegulationByRuleChainId(chainId)
	if gErr != nil {
		return nil, gErr
	}
	return []byte(ruleConfigInfo.RuleConfig), nil
}

// 查询最新修改的一条规则链
func (d *RuleDao) FindLatestComponentRegulation() ([]byte, error) {
	ruleConfigInfo, gErr := FindLatestComponentRegulation()
	if gErr != nil {
		return nil, gErr
	}
	return []byte(ruleConfigInfo.RuleConfig), nil
}

func (d *RuleDao) ListToComponentRegulation(username, keywords string, root *bool, disabled *bool, size, page int) ([]types.RuleChain, int, error) {
	var ruleChains []types.RuleChain
	totalCount := 0
	list, total, err := FindComponentRegulationByPage(page, size, root, disabled, keywords)
	if err != nil {
		return ruleChains, totalCount, err
	}
	totalCount = int(total)
	for _, item := range list {
		var ruleChainItem types.RuleChain
		if err := json.Unmarshal([]byte(item.RuleConfig), &ruleChainItem); err != nil {
			continue
		}
		ruleChains = append(ruleChains, ruleChainItem)
	}
	return ruleChains, totalCount, nil
}

func (d *RuleDao) GetAllComponentRegulation(username string) ([]types.RuleChain, error) {
	var ruleChains []types.RuleChain
	list, err := GetAllLoadComponentRegulation(username)
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		var ruleChainItem types.RuleChain
		if err := json.Unmarshal([]byte(item.RuleConfig), &ruleChainItem); err != nil {
			continue
		}
		ruleChains = append(ruleChains, ruleChainItem)
	}
	return ruleChains, nil
}
