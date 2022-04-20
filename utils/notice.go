package utils

import (
	"encoding/json"
	"fmt"
	models "go-bbs/models/admin"
	"time"
)

type NoticeMessage struct {
	Title     string    `json:"Title"`
	SendId    int       `json:"SendId"`
	ReceiveId int       `json:"ReceiveId"`
	Username  string    `json:"Username"`
	Cover     string    `json:"Cover"`
	UserUrl   string    `json:"UserUrl"`
	Content   string    `json:"Content"`
	Replay    string    `json:"Replay"`
	Url       string    `json:"Url"`
	Date      time.Time `json:"Date"`
}

func SendNotic(m NoticeMessage) (id int64, err error) {
	if m.SendId == m.ReceiveId {
		return
	}
	// fmt.Printf("Notice:%v\n",m)
	jsonContent, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("SendNotic Error : %s\n", err.Error())
		return 0, err
	}
	id, err = models.AddNotice(&models.Notice{
		SendId:    m.SendId,
		ReceiveId: m.ReceiveId,
		Title:     m.Title,
		Content:   string(jsonContent),
		Status:    1,
	})

	if err != nil {
		fmt.Printf("SendNotic Error : %s\n", err.Error())
	}

	return
}
