package model

import "time"

type MdWorkflow struct {
	Id           int64      `gorm:"column:id" json:"id"`
	Title        string     `gorm:"column:title" json:"title"`
	Content      string     `gorm:"column:content" json:"content"`
	Desc         string     `gorm:"column:desc" json:"desc"`
	ChainId      string     `gorm:"column:chain_id" json:"chainId"`
	ChainName    string     `gorm:"column:chain_name" json:"chainName"`
	ChainVersion int64      `gorm:"column:chain_version" json:"chainVersion"`
	CreatedAt    *time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
