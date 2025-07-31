package model

import (
	"time"

	"gorm.io/gorm"
)

type ComponentRegulation struct {
	Id          int64           `gorm:"column:id"`
	UserName    string          `gorm:"column:user_name"`
	Root        bool            `gorm:"column:root"`
	Disabled    bool            `gorm:"column:disabled"`
	RuleChainId string          `gorm:"column:rule_chain_id"`
	Name        string          `gorm:"column:name"`
	RuleConfig  string          `gorm:"column:rule_config"`
	CreatedAt   *time.Time      `gorm:"column:created_at"`
	UpdatedAt   *time.Time      `gorm:"column:updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"column:deleted_at"`
}
