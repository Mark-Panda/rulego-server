package dao

import (
	"fmt"

	"github.com/rulego/rulego-server/internal/model"
)

// 查询所有需要加载的规则链
func GetAllLoadRegulation(username string) ([]model.Regulation, error) {
	re := make([]model.Regulation, 0)
	err := model.DBClient.Client.Model(&model.Regulation{}).Where("disabled = ? and user_name = ?", false, username).Find(&re).Error
	return re, err
}

// 创建规则链
func CreateRegulation(r model.Regulation) error {
	err := model.DBClient.Client.Create(&r).Error
	return err
}

// 根据ID更新规则链
func UpdateRegulationByRuleChainId(ruleChainId string, ruleConfig string) error {
	return model.DBClient.Client.Model(&model.Regulation{}).Where("rule_chain_id = ?", ruleChainId).Update("rule_config", ruleConfig).Error
}

// 根据规则链ID查询规则链信息
func FindRegulationByRuleChainId(ruleChainId string) (*model.Regulation, error) {
	r := &model.Regulation{}
	err := model.DBClient.Client.Model(&model.Regulation{}).Where("rule_chain_id = ?", ruleChainId).Limit(1).Find(&r).Error
	return r, err
}

// 创建规则链
func SaveRegulation(r model.Regulation) error {
	err := model.DBClient.Client.Save(&r).Error
	return err
}

// 根据规则链ID查询规则链信息
func DeleteRegulationByRuleChainId(ruleChainId string) error {
	r := model.Regulation{}
	err := model.DBClient.Client.Model(&model.Regulation{}).Where("rule_chain_id = ?", ruleChainId).Delete(&r).Error
	return err
}

// 查询最新修改的一条规则链
func FindLatestRegulation() (*model.Regulation, error) {
	r := &model.Regulation{}
	err := model.DBClient.Client.Model(&model.Regulation{}).Order("updated_at desc").Limit(1).Find(&r).Error
	return r, err
}

// 分页查询规则链
func FindRegulationByPage(page, size int, root *bool, disabled *bool, keywords string) ([]model.Regulation, int64, error) {
	var re []model.Regulation
	where := ""
	if root != nil {
		if where != "" {
			where = fmt.Sprintf("%s and root = %v  ", where, *root)
		} else {
			where = fmt.Sprintf("root = %v  ", *root)
		}

	}
	if disabled != nil {
		if where != "" {
			where = fmt.Sprintf("%s and disabled = %v  ", where, *disabled)
		} else {
			where = fmt.Sprintf("disabled = %v  ", *disabled)
		}
	}
	if keywords != "" {
		if where != "" {
			where = fmt.Sprintf("%s and name like '%s'  ", where, keywords)
		} else {
			where = fmt.Sprintf("name like '%s'  ", keywords)
		}
	}
	err := model.DBClient.Client.Model(&model.Regulation{}).Where(where).Order("updated_at desc").Offset((page - 1) * size).Limit(size).Find(&re).Error
	total := int64(0)
	model.DBClient.Client.Model(&model.Regulation{}).Where(where).Count(&total)
	return re, total, err
}
