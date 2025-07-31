package dao

import (
	"fmt"

	"github.com/rulego/rulego-server/internal/model"
)

// 创建组件使用规则
func (s *EventDao) CreateComponentUseRule(r model.ComponentUseRule) error {
	err := model.DBClient.Client.Create(&r).Error
	return err
}

// 分页查询组件使用规则
func (s *EventDao) FindComponentUseRuleByPage(page, size int, component_type string, disabled *bool, keywords string) ([]model.ComponentUseRule, int64, error) {
	var re []model.ComponentUseRule
	where := ""
	if disabled != nil {
		if where != "" {
			where = fmt.Sprintf("%s and disabled = %v  ", where, *disabled)
		} else {
			where = fmt.Sprintf("disabled = %v  ", *disabled)
		}
	}
	if keywords != "" {
		if where != "" {
			where = fmt.Sprintf("%s and component_name like '%s'  ", where, "%"+keywords+"%")
		} else {
			where = fmt.Sprintf("component_name like '%s'  ", "%"+keywords+"%")
		}
	}
	if component_type != "" {
		if where != "" {
			where = fmt.Sprintf("%s and component_type = '%s'  ", where, component_type)
		} else {
			where = fmt.Sprintf("component_type = '%s'  ", component_type)
		}
	}
	err := model.DBClient.Client.Model(&model.ComponentUseRule{}).Where(where).Order("updated_at desc").Offset((page - 1) * size).Limit(size).Find(&re).Error
	total := int64(0)
	model.DBClient.Client.Model(&model.ComponentUseRule{}).Where(where).Count(&total)
	return re, total, err
}

// 物理删除组件使用规则
func (s *EventDao) DeleteComponentUseRuleByIdPhysical(id string) error {
	err := model.DBClient.Client.Model(&model.ComponentUseRule{}).Where("id = ?", id).Unscoped().Delete(&model.ComponentUseRule{}).Error
	return err
}

// 根据ID更新组件使用规则
func (s *EventDao) UpdateComponentUseRuleById(id string, data map[string]interface{}) error {
	return model.DBClient.Client.Model(&model.ComponentUseRule{}).Where("id = ?", id).Updates(data).Error
}

// 根据ID查询组件使用规则
func (s *EventDao) FindComponentUseRuleById(id string) (model.ComponentUseRule, error) {
	var r model.ComponentUseRule
	err := model.DBClient.Client.Model(&model.ComponentUseRule{}).Where("id = ?", id).Take(&r).Error
	return r, err
}
