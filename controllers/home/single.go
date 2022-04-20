package home

import (
	"fmt"
	"go-bbs/models/admin"
	"go-bbs/utils"
	"net/http"
	"path"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gocolly/colly"
)

type SingleController struct {
	BaseController
}

func (c *SingleController) About() {
	var nav []map[string]string
	nav = append(nav, map[string]string{
		"Title": "首页",
		"Url":   "/",
	})
	nav = append(nav, map[string]string{
		"Title": "版权申明",
		"Url":   "/about.html",
	})
	c.Data["Nav"] = nav
	c.Data["Navlen"] = len(nav) - 1
	c.TplName = "home/" + c.Template + "/about.html"
}

func (c *SingleController) Sitemap() {

	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, beego.AppPath+"/sitemap.xml")
}

func (c *SingleController) Download() {

	response := make(map[string]interface{})
	customer := c.GetSession("Customer")
	if customer == nil {
		response["msg"] = "请先登录！"
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	url := c.Ctx.Input.Param(":splat")
	downpath := utils.AesDecrypt(url, utils.Key)
	fmt.Printf("url:", downpath)
	fileBaseName := path.Base(downpath)
	c.Ctx.Output.Download(beego.AppPath+downpath, fileBaseName)

}

func (c *SingleController) Links() {

	var nav []map[string]string

	nav = append(nav, map[string]string{
		"Title": "首页",
		"Url":   "/",
	})
	nav = append(nav, map[string]string{
		"Title": "友情链接",
		"Url":   "/links.html",
	})
	c.Data["Nav"] = nav
	c.Data["Navlen"] = len(nav) - 1
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10000
	var offset int64

	sortby = append(sortby, "sort")
	order = append(order, "asc")

	link, _ := admin.GetAllLink(query, fields, sortby, order, offset, limit)
	c.Data["Link"] = link
	c.TplName = "home/" + c.Template + "/links.html"

}

func (c *SingleController) Contact() {
	var nav []map[string]string
	nav = append(nav, map[string]string{
		"Title": "首页",
		"Url":   "/",
	})
	nav = append(nav, map[string]string{
		"Title": "联系我们",
		"Url":   "/contact.html",
	})
	c.Data["Nav"] = nav
	c.Data["Navlen"] = len(nav) - 1
	c.TplName = "home/" + c.Template + "/contact.html"
}

func (c *SingleController) Test() {
	c.TplName = "home/" + c.Template + "/abc.html"
}

func (c *SingleController) Test2() {
	c.TplName = "home/" + c.Template + "/writing.html"
}

func (c *SingleController) Quotes() {
	c.TplName = "home/" + c.Template + "/quotes.html"
}

func (c *SingleController) Olympics() {

	cy := colly.NewCollector()
	html := ""
	cy.OnHTML("#medal-standing table", func(e *colly.HTMLElement) {
		// html,_ = e.DOM.Html()

		e.DOM.Find("thead tr th br").Remove()
		e.DOM.Find("caption").SetText("2021年日本东京奥运会奖牌最新牌行榜 - 实时更新")
		html = utils.OlympicsHtml2UrlAndSrc(e.DOM)
	})
	// On every a element which has f attribute call callback

	cy.Visit("https://olympics.com/tokyo-2020/olympic-games/zh/results/all-sports/medal-standings.htm")

	//html = utils.OlympicsHtml2UrlAndSrc(html)
	c.Data["Html"] = html

	c.TplName = "home/" + beego.AppConfig.String("view") + "/tokyo-2020.html"
}

func (c *SingleController) DailyhotDetail() {
	id := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	var hotnews admin.Hotnews
	qs := o.QueryTable(new(admin.Hotnews))
	err := qs.Filter("id", id).RelatedSel().One(&hotnews)
	if err != nil {
		c.Abort("404")
	}

	c.Data["hotnews"] = hotnews

	var abouthotnews []*admin.Hotnews
	// o = orm.NewOrm()
	o.QueryTable(new(admin.Hotnews)).RelatedSel().OrderBy("-id").Limit(20).All(&abouthotnews)

	c.Data["hotnews"] = hotnews

	for _, v := range abouthotnews {
		// fmt.Printf("Title:%s\n",v.Title)
		v.Url = "/dailyhot/detail/" + strconv.Itoa(v.Id) + ".html"
	}

	c.Data["AboutHotnews"] = abouthotnews

	c.TplName = "home/nihongdengxia/new_detail.html"
}
func (c *SingleController) AgentUrl() {

	strId := c.Ctx.Input.Param(":splat")
	strId = utils.AesDecrypt(strId, utils.Key)
	id, e := strconv.Atoi(strId)
	if e != nil {
		fmt.Println(strId)
		c.Redirect(strId, 302)
	}
	//fmt.Println(id)
	// 基础数据
	o := orm.NewOrm()
	var hotnews admin.Hotnews
	qs := o.QueryTable(new(admin.Hotnews))
	err := qs.Filter("id", id).RelatedSel().One(&hotnews)
	if err != nil {
		c.Abort("404")
	}

	var abouthotnews []*admin.Hotnews
	// o = orm.NewOrm()
	o.QueryTable(new(admin.Hotnews)).RelatedSel().OrderBy("-id").Limit(20).All(&abouthotnews)

	c.Data["hotnews"] = hotnews

	for _, v := range abouthotnews {
		// fmt.Printf("Title:%s\n",v.Title)
		v.Url = "/dailyhot/detail/" + strconv.Itoa(v.Id) + ".html"
	}

	c.Data["AboutHotnews"] = abouthotnews

	c.TplName = "home/nihongdengxia/new_detail.html"

}

func (c *SingleController) History() {
	o := orm.NewOrm()
	var articles []*admin.Article
	qs := o.QueryTable(new(admin.Article))
	qs = qs.Filter("status", 1)
	qs = qs.Filter("Customer__Username__isnull", false)
	qs = qs.Filter("Category__Name__isnull", false)
	// 获取数据
	_, err := qs.OrderBy("-id").All(&articles)
	if err != nil {
		c.Abort("404")
	}
	c.Data["Data"] = &articles
	c.TplName = "home/" + c.Template + "/time.html"
}

func (c *SingleController) SwitchDanmu() {

	status := c.GetString(":splat")
	if status == "off" || status == "on" {
		c.Ctx.SetCookie("SwitchDanmu", status)

	}
	c.Redirect("/", 302)

}

func (c *SingleController) Google17a908a5fab67dec() {
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, beego.AppPath+"/google17a908a5fab67dec.html")
}
