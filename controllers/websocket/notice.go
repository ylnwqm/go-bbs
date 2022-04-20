package websocket

import (
	"encoding/json"
	"time"

	"fmt"
	models "go-bbs/models/admin"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type NoticeController struct {
	beego.Controller
}

func init() {
	go chatroom()
}

var (
	subscribe   = make(chan Subscriber, 10)
	unsubscribe = make(chan int64, 10)
	publish     = make(chan Message, 10)
	clients     = make(map[int64]Subscriber)
)

type Message struct {
	Id      int64  `json:"id"`
	Count   int64  `json:"count"`
	content string `json:"content"`
}

type Subscriber struct {
	Id   int64
	Conn *websocket.Conn
}

func SendMessage(id, count int64, msg string) Message {
	return Message{id, count, msg}
}

func chatroom() {
	for {
		select {
		case sub := <-subscribe:
			if _, ok := clients[sub.Id]; !ok {
				clients[sub.Id] = sub
				// publish <- SendMessage(sub.Id, 0,"连接成功")
			}
			// publish <- SendMessage(sub.Id, 0,"连接成功")

		case message := <-publish:
			broadcastWebSocket(message)

		case unsub := <-unsubscribe:
			if c, ok := clients[unsub]; ok {
				delete(clients, c.Id)
				c.Conn.Close()
				publish <- SendMessage(unsub, 0, "关闭连接")
			}
		}
	}
}

func (c *NoticeController) Join() {
	uid, _ := c.GetInt64("uid")
	t, _ := c.GetInt64("t")
	wskey := fmt.Sprintf("%d%d", uid, t)
	wskeyint, _ := strconv.ParseInt(wskey, 10, 64)

	ws, err := websocket.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil, 1024, 1024)
	if err != nil {
		return
	}

	subscribe <- Subscriber{Id: wskeyint, Conn: ws}

	//发生异常时，通知关闭该用户
	defer func() {
		unsubscribe <- wskeyint
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// notice
			noticeCount, _ := models.GetNoticeCount(int(uid))
			publish <- SendMessage(wskeyint, noticeCount, "获取消息数量成功")
		}
	}
}

func broadcastWebSocket(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}

	if c, ok := clients[msg.Id]; ok {
		if c.Conn != nil {
			if c.Conn.WriteMessage(websocket.TextMessage, data) != nil {
				unsubscribe <- msg.Id //通知关闭该用户
			}
		}
	}
}
