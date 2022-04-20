package hotnews

import (
	"fmt"
	"go-bbs/models/admin"
	"time"

	"github.com/gocolly/colly"
	// "go-bbs/utils"
	// "strconv"
)

type Weibo struct {
}

func (w Weibo) Get(url string) {

	c := colly.NewCollector()
	c.OnHTML("#pl_top_realtimehot table tbody tr", func(e *colly.HTMLElement) {

		url := e.ChildAttr(".td-02 a", "href")
		// fmt.Println(url)
		title := e.ChildText(".td-02 a")
		url = "https://s.weibo.com" + url
		_, err := admin.AddHotnews(&admin.Hotnews{
			Category: &admin.Category{Id: 18},
			Title:    title,
			Url:      url,
			Created:  time.Now(),
		})

		if err != nil {
			fmt.Println("Weibo Error : " + err.Error())
		} else {
			// utils.LPush(QueryName,strconv.FormatInt(id,10) + "," + url)
		}
	})
	c.Visit(url)
}
