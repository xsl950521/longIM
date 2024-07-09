package models

import (
	"LongIM/utils"
	"gorm.io/gorm"
)

// 人员关系

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int // 关系类型 0-好友
	Desc     string
}

func (m *Contact) TableName() string {
	result := utils.DB.Migrator().HasTable("contact")
	if !result {
		utils.DB.AutoMigrate(&Message{})
	}
	return "contact"
}
