package model

import (
	"time"
)

type RunLog struct {
	Id             int64      `gorm:"column:id"`
	RunId          string     `gorm:"column:run_id"`
	ChainId        string     `gorm:"column:chain_id"`
	ChainName      string     `gorm:"column:chain_name"`
	NodeLog        string     `gorm:"column:node_log"`
	AdditionalInfo string     `gorm:"column:additional_info"`
	RuleChainInfo  string     `gorm:"column:rule_chain_info"`
	Metadata       string     `gorm:"column:metadata"`
	StartTs        int64      `gorm:"column:start_ts"`
	EndTs          int64      `gorm:"column:end_ts"`
	CreatedAt      *time.Time `gorm:"column:created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at"`
}
