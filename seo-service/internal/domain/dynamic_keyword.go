package domain

import (
	"time"
)

type DynamicKeyword struct {
	Id          int64     `gorm:"column:id;primaryKey" json:"id"`
	Url         string    `gorm:"column:url;type:text" json:"url"`
	Description string    `gorm:"column:description;type:text"  json:"description"`
	IsActive    bool      `gorm:"column:is_active"  json:"is_active"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (u *DynamicKeyword) TableName() string {
	return "dynamic_keyword"
}
