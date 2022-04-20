package hotnews

import (
	"fmt"
	"go-bbs/models/admin"
	"time"

	"github.com/gocolly/colly"
	//"go-bbs/utils"
	//"strconv"
)

type Hupu struct {
}

func (h Hupu) Get(url string) {

	c := colly.NewCollector()
	c.OnHTML(".infinite-container .list-item-info", func(e *colly.HTMLElement) {
		url := e.ChildAttr(".list-item-title", "href")
		title := e.ChildText(".item-title-conent")
		_, err := admin.AddHotnews(&admin.Hotnews{
			Category: &admin.Category{Id: 17},
			Title:    title,
			Url:      url,
			Created:  time.Now(),
		})

		if err != nil {
			fmt.Println(err.Error())
		} else {

			// utils.LPush(QueryName,strconv.FormatInt(id,10)+ "," + url)
		}

	})
	c.Visit(url)
}

func (h Hupu) GetMatch(url string) []map[string]string {
	c := colly.NewCollector()
	var ss []map[string]string
	c.OnHTML(".cardListContainer .slick-list .slick-track .slick-slide", func(e *colly.HTMLElement) {

		var s []string
		e.ForEach(".cardItem div .cardItemTitle span", func(_ int, el *colly.HTMLElement) {
			s = append(s, el.Text)
		})

		e.ForEach(".cardItem div .cardItemContent .cardItemContentLine div", func(_ int, el *colly.HTMLElement) {
			if el.Attr("class") == "cardImg" {
				s = append(s, el.ChildAttr(".cardImg img", "src"))
				s = append(s, el.Text)
			} else {
				if el.ChildAttr("div", "class") == "cardItemContentLine-count" {
					var sc string

					sc = el.ChildText("div div div")
					if sc == "" {
						sc = el.ChildText("div div")
					}
					s = append(s, sc)
				} else {
					var sc string
					sc = el.ChildText("div div")
					if sc == "" {
						sc = el.ChildText("div")
					}

					s = append(s, sc)
				}
			}
		})

		// fmt.Printf("%v\n",s)
		var sv []string
		for _, v := range s {
			if v != "" {
				sv = append(sv, v)
			}
		}

		if len(sv) >= 8 {
			ss = append(ss, map[string]string{
				"type":       sv[0],
				"status":     sv[1],
				"team1img":   sv[2],
				"team1name":  sv[3],
				"team1score": sv[4],
				"team2img":   sv[5],
				"team2name":  sv[6],
				"team2score": sv[7],
			})
		}

	})
	c.Visit(url)
	for i, j := 0, len(ss)-1; i < j; i, j = i+1, j-1 {
		ss[i], ss[j] = ss[j], ss[i]
	}
	return ss
}
