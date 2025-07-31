package model

import (
	"time"
)

type ComponentUseRule struct {
	Id            int64      `gorm:"column:id;primaryKey" json:"id"`
	ComponentName string     `gorm:"column:component_name" json:"componentName"`
	ComponentType string     `gorm:"column:component_type" json:"componentType"`
	Disabled      bool       `gorm:"column:disabled" json:"disabled"`
	UseDesc       string     `gorm:"column:use_desc" json:"useDesc"`
	UseRuleDesc   string     `gorm:"column:use_rule_desc" json:"useRuleDesc"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
