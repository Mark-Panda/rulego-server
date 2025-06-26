package service

import (
	"github.com/rulego/rulego-server/config"
	"github.com/rulego/rulego-server/internal/model"
)

func Setup(config config.Config) error {

	if err := model.StartDB(config); err != nil {
		return err
	}

	if s, err := NewUserService(config); err != nil {
		return err
	} else {
		UserServiceImpl = s
	}
	if s, err := NewUserRuleEngineServiceImpl(config); err != nil {
		return err
	} else {
		UserRuleEngineServiceImpl = s
	}

	if s, err := NewEventService(config); err != nil {
		return err
	} else {
		EventServiceImpl = s
	}
	return nil
}
