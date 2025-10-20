package service

import "github.com/rulego/rulego-server/internal/model"

// 创建Md工作流
func (s *EventService) CreateMdWorkflow(info model.MdWorkflow) error {
	return s.EventDao.SaveMdWorkFlow(info)
}

// 编辑Md工作流
func (s *EventService) EditMdWorkflow() error {
	return s.EventDao.EditorMdWorkFlow()
}

// 分页获取md工作流列表
func (s *EventService) GetMdWorkflowList(title string, current, size int) ([]*model.MdWorkflow, int64, error) {
	return s.EventDao.GetMdWorkflowList(title, current, size)
}
