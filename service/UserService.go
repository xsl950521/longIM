package service

import (
	"LongIM/models"
	"LongIM/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
)

// GetUserList
// @Summary 获取用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList  [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}

// FindUserByNameAndPwd
// @Summary 登陆
// @Tags 用户模块
// @param name query string false "用户名"
// @param pass_word query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd  [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	name := c.Query("name")
	pwd := c.Query("pass_word") //输入密码
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(-1, gin.H{
			"message": "该用户不存在",
		})
		return
	}
	fmt.Println("user=", user)
	if !utils.ValidPassword(pwd, user.Salt, user.PassWord) {
		c.JSON(-1, gin.H{
			"message": "密码不正确",
		})
		return
	}
	pwd = utils.MakePassword(pwd, user.Salt)
	fmt.Println("pwd=", pwd)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser  [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "用户名已注册",
		})
		return
	}
	//TODO 手机号、邮箱的重复验证
	if password != repassword {
		c.JSON(-1, gin.H{
			"message": "两次密码不一致",
		})
		return
	}
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	//user.PassWord = password
	db := models.CreateUser(user)
	if db.Error != nil {
		c.JSON(-1, gin.H{
			"message": db.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "新增用户成功",
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser  [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	ID, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(ID)
	db := models.DeleteUser(user)
	if db.Error != nil {
		c.JSON(-1, gin.H{
			"message": db.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除用户成功",
	})
}

// UpdateUser
// @Summary 更新用户
// @Tags 用户模块
// @param id formData string false "id"
// @param password formData string false "password"
// @param name formData string false "name"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser  [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	ID, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(ID)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
		return
	}
	db := models.UpdateUser(user)
	if db.Error != nil {
		c.JSON(-1, gin.H{
			"message": db.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}
