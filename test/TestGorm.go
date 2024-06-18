package test

import (
	"LongIM/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func TestGorm() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&models.UserBasic{})
	if err != nil {
		return
	}

	user := &models.UserBasic{}
	user.Name = "许士龙"
	user.LoginTime = time.Now()
	user.LogoutTime = time.Now()
	user.HeartBeatTime = time.Now()
	db.Create(user)

	fmt.Println(db.First(user, 1))

	db.Model(user).Update("PassWord", "123456")
}
