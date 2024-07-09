package router

import (
	"LongIM/docs"
	"LongIM/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/index", service.GetIndex)

	// 静态资源
	r.Static("/assert", "assert/")
	r.LoadHTMLGlob("views/**/*")
	// 首页
	r.GET("/index", service.GetIndex)
	r.GET("/user/getUserList", service.GetUserList)
	r.GET("/user/createUser", service.CreateUser)
	r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)

	// 发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendPrivateMsg", service.SendPrivateMsg)
	return r
}
