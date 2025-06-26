package dao

import (
	"github.com/dromara/carbon/v2"
	"github.com/rulego/rulego-server/config"
	"github.com/rulego/rulego-server/internal/model"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/json"
)

func NewRuleDaoToDataBase(config config.Config, username string) (*RuleDao, error) {
	dao := &RuleDao{
		config:   config,
		username: username,
		index:    Index{Rules: make(map[string]RuleMeta)},
	}

	return dao, nil
}

// 保存或更新到数据库
func (d *RuleDao) SaveToDataBase(username, chainId string, def []byte) error {
	v, _ := json.Format(def)
	ruleConfigInfo, gErr := FindRegulationByRuleChainId(chainId)
	if gErr != nil {
		return gErr
	}
	if ruleConfigInfo != nil && ruleConfigInfo.RuleChainId != "" {
		return UpdateRegulationByRuleChainId(chainId, string(v))
	}
	// def 转成 types.RuleChain
	var ruleChain types.RuleChain
	if err := json.Unmarshal(def, &ruleChain); err != nil {
		return err
	}
	t := carbon.Now(carbon.Shanghai).StdTime()
	createInfo := model.Regulation{
		UserName:    username,
		Root:        ruleChain.RuleChain.Root,
		Disabled:    ruleChain.RuleChain.Disabled,
		Name:        ruleChain.RuleChain.Name,
		RuleChainId: chainId,
		RuleConfig:  string(v),
		CreatedAt:   &t,
		UpdatedAt:   &t,
	}
	return CreateRegulation(createInfo)
}

// 从数据库删除规则链
func (d *RuleDao) DeleteToDataBase(username, chainId string) error {
	return DeleteRegulationByRuleChainId(chainId)
}

// 按规则链id从数据库查询规则链
func (d *RuleDao) FindDataBaseByRuleChainId(chainId string) ([]byte, error) {
	ruleConfigInfo, gErr := FindRegulationByRuleChainId(chainId)
	if gErr != nil {
		return nil, gErr
	}
	return []byte(ruleConfigInfo.RuleConfig), nil
}

// 查询最新修改的一条规则链
func (d *RuleDao) FindLatestDataBase() ([]byte, error) {
	ruleConfigInfo, gErr := FindLatestRegulation()
	if gErr != nil {
		return nil, gErr
	}
	return []byte(ruleConfigInfo.RuleConfig), nil
}

func (d *RuleDao) ListToDataBase(username, keywords string, root *bool, disabled *bool, size, page int) ([]types.RuleChain, int, error) {
	var ruleChains []types.RuleChain
	totalCount := 0
	list, total, err := FindRegulationByPage(page, size, root, disabled, keywords)
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
