package utils

import (
    "io/ioutil"
    "net/http"
	"encoding/json"
)

type DailyContent struct{
	ImgUrl string `json:"img_url"`
	Forward string `json:"forward"`
	WordsInfo string `json:"words_info"`
}

type DailyData struct{
	DailyContent []DailyContent `json:"content_list"`
}

type DailyRet struct{
	Data DailyData `json:"data"`
}

type Hitokoto struct{
	Content string `json:"hitokoto"`
	From string `json:"from"`
	Creator string `json:"creator"`
}

func GetDailyData()(dailyInfo *DailyRet,err error){

    resp, err := http.Get("http://v3.wufazhuce.com:8000/api/channel/one/0/0")
    if err != nil {
		return nil,err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	json.Unmarshal(body,&dailyInfo)

	if len(dailyInfo.Data.DailyContent) > 0 {
		dailyInfo.Data.DailyContent = dailyInfo.Data.DailyContent[:1]
		imgurl ,err := DownImage(dailyInfo.Data.DailyContent[0].ImgUrl)
		if err == nil{
			dailyInfo.Data.DailyContent[0].ImgUrl = imgurl
		}
	}

	resp, err = http.Get("https://v1.hitokoto.cn/?c=d")
    if err != nil {
		return nil,err
    }
    defer resp.Body.Close()

    body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	var h Hitokoto
	json.Unmarshal(body,&h)
	dailyInfo.Data.DailyContent[0].Forward = h.Content
	dailyInfo.Data.DailyContent[0].WordsInfo = h.From
	return 
}
