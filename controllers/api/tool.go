package api

import (
	// "go-bbs/models/admin"
	// "net/http"
	// "encoding/json"
	// "io/ioutil"
	// "fmt"
	// "time"
	"go-bbs/utils/sitemap"
	// "go-bbs/utils"
)

type ToolController struct {
	BaseController
}

type D struct {
	CreateTime int    `json:"CreateTime"`
	Title      string `json:"Title"`
	Url        string `json:"Url"`
	ImgUrl     string `json:"imgUrl"`
}

type dd struct {
	DH []D `json:"data"`
}

type H struct {
	Data dd `json:"Data"`
}

func (ctl *ToolController) CSitemap() {
	// d,err := utils.GetDailyData()
	//fmt.Printf("%v\n",d)
	//fmt.Printf("%v\n",err)
	sitemap.Sitemap("./views/home/toutiao", "http://clblog.club")
	// ctl.Data["json"] = d
	// ctl.ServeJSON()
	// ctl.StopRun()

	//var w = utils.Williamlong{}
	//for i:=2;i<=13;i++ {
	//w.Get("https://www.williamlong.info/book/page/index_"+ strconv.Itoa(i)+".html")
	//}
	//w.Get("https://www.williamlong.info/book/page/index_1.html")
	//ctl.ServeJSON()
	//ctl.StopRun()

	// resp, _ := http.Get("https://api.tophub.fun/v2/GetAllInfoGzip?id=1065&page=0&type=pc")
	// defer resp.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(body)
	// var r H
	// err := json.Unmarshal([]byte(body), &r)

	// if err != nil{
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(r)

	// for _,v := range r.Data.DH {

	// 	tm := time.Unix(int64(v.CreateTime), 0)
	// 	admin.AddHotnews(&admin.Hotnews{
	// 		Category:&admin.Category{Id:1},
	// 		Title:v.Title,
	// 		Url:v.Url,
	// 		Created:tm,
	// 	})
	// }

}
