package service

import (
	"github.com/gin-gonic/gin"
	"text/template"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index  [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
		return
	}
	ind.Execute(c.Writer, "index")
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "welcome !!",
	//})
}
