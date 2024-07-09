package models

import (
	"LongIM/utils"
	"gorm.io/gorm"
)

type GroupBasic struct {
	gorm.Model
	Name    string `gorm:"type:varchar(20);not null"`
	OwnerId uint   `gorm:"type:int;not null"`
	Icon    string `gorm:"type:varchar(255);"`
	Type    int    `gorm:"type:int;"`
	Level   int    `gorm:"type:int;not null"`
}

func (m *GroupBasic) TableName() string {
	result := utils.DB.Migrator().HasTable("group")
	if !result {
		utils.DB.AutoMigrate(&Message{})
	}
	return "group"
}
