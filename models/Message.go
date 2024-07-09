package models

import (
	"LongIM/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	FromId   int64  // 发送者
	TargetId int64  // 接受者
	Type     uint   // 消息类型 私聊。群聊，广播
	Media    int    // 消息类型 文字 图片 音频
	Content  string // 消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int // 其他数字统计
}

func (m *Message) TableName() string {
	result := utils.DB.Migrator().HasTable("message")
	if !result {
		utils.DB.AutoMigrate(&Message{})
	}
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQue   chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLock sync.RWMutex

func Chat(w http.ResponseWriter, request *http.Request) {
	// TODO：校验token等合法性
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	//token := query.Get("token")

	tid := query.Get("targetId")
	targetId, _ := strconv.ParseInt(tid, 10, 64)
	//context := query.Get("context")
	//sendType := query.Get("sendType")
	isValida := true
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(w, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取连接
	node := &Node{
		Conn:      conn,
		DataQue:   make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	// TODO:用户关系
	//userid 跟 node绑定并加锁
	rwLock.Lock()
	clientMap[userId] = node
	rwLock.Unlock()
	// 完成发送逻辑
	go sendProc(node)
	// 完成接受逻辑
	go recProc(node)
	sendMsg(userId, targetId, []byte("欢迎进入聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws] <<<<<< ", string(data))
	}
}

var udpSendChan = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpSendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecProc()
}

// 完成udp 数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case data := <-udpSendChan:
			_, err = con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
func udpRecProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		var buf [1024]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私聊
		sendMsg(msg.FromId, msg.TargetId, data)
		//case 2: // 群聊
		//	sendGroupMsg()
		//case 3: // 广播
		//	sendAllMsg()
		//case 4:
	}

}

func sendMsg(userId int64, targetId int64, msg []byte) {
	rwLock.RLock()
	node, ok := clientMap[userId]
	rwLock.RUnlock()
	if ok {
		node.DataQue <- msg
	}

}
