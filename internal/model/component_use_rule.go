package model

import (
	"time"
)

type ComponentUseRule struct {
	Id            int64      `gorm:"column:id"`
	ComponentName string     `gorm:"column:component_name"`
	ComponentType string     `gorm:"column:component_type"`
	Disabled      bool       `gorm:"column:disabled"`
	UseDesc       string     `gorm:"column:use_desc"`
	UseRuleDesc   string     `gorm:"column:use_rule_desc"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
}
