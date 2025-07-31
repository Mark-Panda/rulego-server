package dao

import (
	"fmt"

	"github.com/rulego/rulego-server/internal/model"
)

// 查询所有需要加载的规则链
func GetAllLoadComponentRegulation(username string) ([]model.ComponentRegulation, error) {
	re := make([]model.ComponentRegulation, 0)
	err := model.DBClient.Client.Model(&model.ComponentRegulation{}).Where("disabled = ? and user_name = ?", false, username).Find(&re).Error
	return re, err
}

// 创建规则链
func CreateComponentRegulation(r model.ComponentRegulation) error {
	err := model.DBClient.Client.Create(&r).Error
	return err
}

// 根据ID更新规则链
func UpdateComponentRegulationByRuleChainId(ruleChainId string, data map[string]interface{}) error {

	return model.DBClient.Client.Model(&model.ComponentRegulation{}).Where("rule_chain_id = ?", ruleChainId).Updates(data).Error
}

// 根据规则链ID查询规则链信息
func FindComponentRegulationByRuleChainId(ruleChainId string) (*model.ComponentRegulation, error) {
	r := &model.ComponentRegulation{}
	err := model.DBClient.Client.Model(&model.ComponentRegulation{}).Where("rule_chain_id = ?", ruleChainId).Limit(1).Find(&r).Error
	return r, err
}

// 创建规则链
func SaveComponentRegulation(r model.ComponentRegulation) error {
	err := model.DBClient.Client.Save(&r).Error
	return err
}

// 根据规则链ID查询规则链信息
func DeleteComponentRegulationByRuleChainId(ruleChainId string) error {
	r := model.ComponentRegulation{}
	err := model.DBClient.Client.Model(&model.ComponentRegulation{}).Where("rule_chain_id = ?", ruleChainId).Delete(&r).Error
	return err
}

// 物理删除规则链
func DeleteComponentRegulationByRuleChainIdPhysical(ruleChainId string) error {
	err := model.DBClient.Client.Model(&model.ComponentRegulation{}).Where("rule_chain_id = ?", ruleChainId).Unscoped().Delete(&model.Regulation{}).Error
	return err
}

// 查询最新修改的一条规则链
func FindLatestComponentRegulation() (*model.ComponentRegulation, error) {
	r := &model.ComponentRegulation{}
	err := model.DBClient.Client.Model(&model.ComponentRegulation{}).Order("updated_at desc").Limit(1).Find(&r).Error
	return r, err
}

// 分页查询规则链
func FindComponentRegulationByPage(page, size int, root *bool, disabled *bool, keywords string) ([]model.ComponentRegulation, int64, error) {
	var re []model.ComponentRegulation
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
			where = fmt.Sprintf("%s and name like '%s'  ", where, "%"+keywords+"%")
		} else {
			where = fmt.Sprintf("name like '%s'  ", "%"+keywords+"%")
		}
	}
	err := model.DBClient.Client.Model(&model.ComponentRegulation{}).Where(where).Order("updated_at desc").Offset((page - 1) * size).Limit(size).Find(&re).Error
	total := int64(0)
	model.DBClient.Client.Model(&model.ComponentRegulation{}).Where(where).Count(&total)
	return re, total, err
}
