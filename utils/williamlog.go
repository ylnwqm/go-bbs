package utils

import (
	"fmt"
	"go-bbs/models/admin"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/gocolly/colly"
)

type Williamlong struct {
}

func (w Williamlong) Get(url string) {

	//var sli []map[string]string
	c := colly.NewCollector()
	c.OnHTML("#main .content .entry", func(e *colly.HTMLElement) {

		//s := make(map[string]string)
		href := e.ChildAttr("h1 .post-title", "href")
		if href != "" {
			fmt.Println(href)
			GetDetail(href)
		}
		//s["title"] = e.ChildText("h1 .post-title")
		//fmt.Println(s)
	})

	c.Visit(url)
}

func GetDetail(url string) {
	c := colly.NewCollector()
	c.OnHTML("#main .content .entry", func(e *colly.HTMLElement) {
		s := make(map[string]string)
		s["title"] = e.ChildText("h1 .post-title")
		author := e.ChildText(".post-meta")
		str := strings.Split(author, "|")
		for _, v := range str {
			if v != "" {
				info := strings.Split(v, ":")
				if info[0] == "作者" {
					s["author"] = info[1]
					break
				}
			}
		}
		//s["author"] = e.ChildText(".post-meta")
		s["content"] = e.ChildText(".entry-content")
		//fmt.Println(s["title"])
		save(s)
	})

	c.Visit(url)
}

var worker, _ = NewIdWorker(1)

func save(a map[string]string) {

	if len(a["content"]) == 0 || len(a["author"]) == 0 || len(a["title"]) == 0 {
		return
	}
	fmt.Println(a["title"])
	html := HtmlRemoteImg2LocalImg(a["content"])
	md := Html2md(a["content"])

	uid, err := worker.GetNextId()
	if err != nil {
		uid, err = worker.GetNextId()
		if err != nil {
			uid = time.Now().Unix()
		}
	}
	uidstr := strconv.FormatInt(uid, 10)

	customerId, err := admin.AddCustomer(&admin.Customer{
		Uid:      uidstr,
		Username: a["author"],
		Nickname: a["author"],
		Image:    "/static/images/cover_default.jpg",
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	o := orm.NewOrm()
	art := admin.Article{
		Cover:    "/static/images/article_default.jpg",
		Review:   0,
		Title:    a["title"],
		Remark:   Subbbs(a["content"], 200, "javascript:;"),
		Desc:     md,
		Html:     html,
		Category: &admin.Category{Id: 13},
		Customer: &admin.Customer{Id: int(customerId)},
		Status:   1,
		Other:    "",
		Tag:      a["title"] + "," + a["author"],
	}

	id, err := o.Insert(&art)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("新增成功：ID : %d,标题：%s\n", id, a["title"])
	}

}
