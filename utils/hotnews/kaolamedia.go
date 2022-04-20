package hotnews

import (
	"fmt"
	"go-bbs/models/admin"
	"time"

	"github.com/gocolly/colly"
	//"go-bbs/utils"
	//"strconv"
)

type Kaolamedia struct {
}

func (k Kaolamedia) Get(url string) {

	c := colly.NewCollector()
	c.OnHTML(".container .uk-grid div .hot-block", func(e *colly.HTMLElement) {

		ty := e.ChildText(".title")
		//fmt.Println(ty)
		//fmt.Println(ty)
		//e.DOM.Children().Filter(".list .head").Remove()
		e.ForEach(".list .row", func(_ int, el *colly.HTMLElement) {
			url := el.ChildAttr(".keyword", "href")
			title := el.ChildText(".keyword")
			hDesc := el.ChildText(".value")

			if url != "" {
				hw := &admin.Hotnews{
					//Category  *Category `orm:"rel(one)"`
					Title:   title,
					Url:     url,
					HotDesc: hDesc,
					Created: time.Now(),
				}
				switch ty {
				case "百度热点":
					hw.Category = &admin.Category{Id: 27}
				case "微博热点":
					hw.Category = &admin.Category{Id: 28}
				case "知乎热榜":
					hw.Category = &admin.Category{Id: 21}
				case "百度贴吧":
					hw.Category = &admin.Category{Id: 25}
				case "哔哩哔哩":
					hw.Category = &admin.Category{Id: 20}
				case "抖音热榜":
					hw.Category = &admin.Category{Id: 22}
				case "抖音热搜":
					hw.Category = &admin.Category{Id: 26}
				case "豆瓣话题":
					hw.Category = &admin.Category{Id: 24}
				case "网易新闻":
					hw.Category = &admin.Category{Id: 23}
				}

				_, err := admin.AddHotnews(hw)
				if err != nil {
					fmt.Println("Kaolamedia Error : " + err.Error())
				} else {
					//utils.LPush(QueryName,strconv.FormatInt(id,10) + "," + url)
				}
			}
		})
		// fmt.Println(title)
		// _,err := admin.AddHotnews(&admin.Hotnews{
		// 		Category:&admin.Category{Id:19},
		// 		Title:title,
		// 		Url:url,
		// 		Created:time.Now(),
		// 	})

		// if err != nil {
		// 	fmt.Println("Baidu Error : "+err.Error())
		// }
	})
	c.Visit(url)
}
