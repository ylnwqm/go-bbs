package hotnews
import (
	"github.com/gocolly/colly"
)

type Page struct{

}

func (p Page) Get(url string) string {

	var html string
	c := colly.NewCollector()
	c.OnHTML(".main-thread .thread-content-detail", func(e *colly.HTMLElement) {
		html,_ = e.DOM.Html()
	})
	c.Visit(url)

	return html
}