package hotnews

import (
	"fmt"
	"go-bbs/models/admin"
	"time"

	"github.com/gocolly/colly"
	// "go-bbs/utils"
	// "strconv"
)

type Baidu struct {
}

func (b Baidu) Get(url string) {

	c := colly.NewCollector()
	c.OnHTML(".right-container_2EFJr .container-bg_lQ801 .horizontal_1eKyQ .content_1YWBm", func(e *colly.HTMLElement) {

		url := e.ChildAttr("a", "href")
		title := e.ChildText(".title_dIF3B")
		// fmt.Println(title)
		_, err := admin.AddHotnews(&admin.Hotnews{
			Category: &admin.Category{Id: 19},
			Title:    title,
			Url:      url,
			Created:  time.Now(),
		})

		if err != nil {
			fmt.Println("Baidu Error : " + err.Error())
		} else {

			// utils.LPush(QueryName,strconv.FormatInt(id,10) + "," + url)
		}
	})
	c.Visit(url)
}
