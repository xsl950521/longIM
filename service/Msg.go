package service

import (
	"LongIM/models"
	"LongIM/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// 防止跨域
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println("err:", err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println("err:", err)
		}
		fmt.Println("msg:", msg)
		tNow := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tNow, msg)
		fmt.Println("发送消息：", m)
		err = ws.WriteMessage(websocket.TextMessage, []byte(m))
		if err != nil {
			fmt.Println("err:", err)
		}
	}

}

func SendPrivateMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
