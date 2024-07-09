package models

import (
	"LongIM/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	ClientIP      string
	Identity      string
	LoginTime     time.Time `gorm:"type:datetime;default:null"`
	HeartBeatTime time.Time `gorm:"type:datetime;default:null"`
	LogoutTime    time.Time `gorm:"type:datetime;default:null"`
	Salt          string
}

func (t *UserBasic) TableName() string {
	return "user_basic"
}

func FindUserByNameAndPwd(name, pwd string) UserBasic {
	user := UserBasic{}
	err := utils.DB.Where("name = ? and pass_word = ?", name, pwd).First(&user).Error
	if err != nil {
		panic(err)
	}
	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.Md5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func CreateUser(user UserBasic) *gorm.DB {
	err := utils.DB.AutoMigrate(&UserBasic{})
	if err != nil {
		return nil
	}
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{
		Name:      user.Name,
		PassWord:  user.PassWord,
		LoginTime: time.Now(),
		Phone:     user.Phone,
		Email:     user.Email,
	})
}
