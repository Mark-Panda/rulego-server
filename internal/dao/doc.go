package dao

import "github.com/rulego/rulego-server/internal/model"

func (s *EventDao) SaveMdWorkFlow(mdWorkflow model.MdWorkflow) error {
	return model.DBClient.Client.Create(&mdWorkflow).Error
}

func (s *EventDao) EditorMdWorkFlow(mdId int64, updateInfo map[string]interface{}) error {
	return model.DBClient.Client.Model(&model.MdWorkflow{}).Where("id = ?", mdId).Updates(updateInfo).Error
}

func (s *EventDao) GetMdWorkflowList(title string, current, size int) ([]*model.MdWorkflow, int64, error) {
	var total int64
	var result []*model.MdWorkflow
	where := ""
	if title != "" {
		where = "title like '%" + title + "%'"
	}
	if err := model.DBClient.Client.Where(where).Order("created_at desc").Offset((current - 1) * size).Limit(size).Find(&result).Error; err != nil {
		return result, 0, err
	}
	if err := model.DBClient.Client.Model(&model.MdWorkflow{}).Where(where).Count(&total).Error; err != nil {
		return result, 0, err
	}
	return result, total, nil
}

// 删除md工作流
func (s *EventDao) DeleteMdWorkflow(id int) error {
	return model.DBClient.Client.Delete(&model.MdWorkflow{}, id).Error
}
