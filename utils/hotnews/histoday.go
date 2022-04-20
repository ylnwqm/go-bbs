package hotnews
import (
	"github.com/gocolly/colly"
	"strings"
)

type Histoday struct{

}

func (h Histoday) Get(url string) ([]map[string]string){

	c := colly.NewCollector()
	var res []map[string]string
	c.OnHTML(".tih-list > .tih-item", func(e *colly.HTMLElement) {

		var temp = make(map[string]string)
		title := e.ChildText("dt")
		title = strings.TrimLeft(strings.Split(title, ".")[1]," ")
		img := e.ChildAttr("dd .item-img","data-src")
		if img == "" {
			img = e.ChildAttr("dd .item-img","src")
		}
		url := e.ChildAttr(".clearfix .right-btn-container a","href")
		// fmt.Printf("IMG:%s\n",img)
		temp["title"] = title
		temp["img"] = img
		temp["url"] = url
		res = append(res,temp)
	})
	c.Visit(url)
	return res
}
