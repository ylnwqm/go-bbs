package home

import (
	"go-bbs/models/admin"
	"go-bbs/utils"
	"os"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type HotnewsController struct {
	BaseController
}

func (c *HotnewsController) Get() {
	response := make(map[string]interface{})

	limit := int64(30)
	page, _ := c.GetInt64("page", 1) // 页数
	offset := (page - 1) * limit     // 偏移量

	var hotnews []*admin.Hotnews
	o := orm.NewOrm()
	o.QueryTable(new(admin.Hotnews)).RelatedSel().OrderBy("-id").Limit(limit).Offset(offset).All(&hotnews)

	for _, v := range hotnews {
		//fmt.Println(v.Category)
		//v.Created = v.Created
		v.Rurl = v.Url
		v.Url = "/dailyhot/detail/" + strconv.Itoa(v.Id) + ".html"
	}

	count, _ := o.QueryTable(new(admin.Hotnews)).RelatedSel().Count()

	response["Data"] = hotnews
	response["Paginator"] = utils.GenPaginator(page, limit, count)

	response["msg"] = "Success."
	response["code"] = 200
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *HotnewsController) List() {
	o := orm.NewOrm()
	var hotnews []*admin.Hotnews
	o.QueryTable(new(admin.Hotnews)).RelatedSel().OrderBy("-id").Limit(30).All(&hotnews)

	for _, v := range hotnews {
		//fmt.Println(v.Category)
		v.Rurl = v.Url
		if c.Template == "baidu" {
			v.Url = "/hotnews/detail/" + strconv.Itoa(v.Id) + ".html"
		} else {
			v.Url = "/dailyhot/detail/" + strconv.Itoa(v.Id) + ".html"
		}
	}

	c.Data["Hotnews"] = hotnews

	c.TplName = "home/" + c.Template + "/hot.html"

}

func (c *HotnewsController) GetInc() {
	response := make(map[string]interface{})

	var ret []map[string]string

	var hotnews []*admin.Hotnews
	o := orm.NewOrm()
	o.QueryTable(new(admin.Hotnews)).RelatedSel().Limit(50).OrderBy("-id").All(&hotnews)

	for _, v := range hotnews {
		var item = make(map[string]string)
		item["href"] = "/hotnews/detail/" + strconv.Itoa(v.Id) + ".html"

		var r = []rune(v.Title)
		length := len(r)
		if length < 26 {
			item["info"] = v.Title
		} else {
			item["info"] = string(r[0:26]) + ` ...`
		}
		item["img"] = "/static/danmu/static/img/heisenberg.png"

		ret = append(ret, item)
	}

	response["Data"] = ret
	response["msg"] = "Success."
	response["code"] = 200
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *HotnewsController) GetHistoryHot() {

	var hotnews []*admin.Hotnews
	o := orm.NewOrm()
	o.QueryTable(new(admin.Hotnews)).RelatedSel().OrderBy("-Created").All(&hotnews, "Created")

	var datetime = make(map[string]int64)
	var dateTimeKey []string
	for _, v := range hotnews {
		k := v.Created.Format("2006-01-02")
		if datetime[k] == 0 {
			dateTimeKey = append(dateTimeKey, k)
		}
		datetime[k] = datetime[k] + 1
	}

	var hotnewsData []*admin.Hotnews
	o.QueryTable(new(admin.Hotnews)).RelatedSel().Limit(30).OrderBy("-Created").All(&hotnewsData)
	for _, v := range hotnewsData {
		//fmt.Println(v.Category)
		v.Rurl = v.Url
		v.Url = "/hotnews/detail/" + strconv.Itoa(v.Id) + ".html"
	}

	c.Data["Data"] = hotnewsData
	c.Data["DateTime"] = datetime
	c.Data["DateTimeKey"] = dateTimeKey
	c.TplName = "home/" + c.Template + "/more.html"
}

func (c *HotnewsController) GetHistoryHotByDate() {
	response := make(map[string]interface{})

	var hotnews []*admin.Hotnews
	o := orm.NewOrm()
	date := c.GetString("date")
	if date == "" || len(date) != 10 {
		response["Data"] = "非法操作！"
		response["msg"] = "Success."
		response["code"] = 200
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	o.QueryTable(new(admin.Hotnews)).Filter("created__gte", date+" 00:00:00").Filter("created__lte", date+" 23:59:59").RelatedSel().Limit(20).OrderBy("-Created").All(&hotnews)

	for _, v := range hotnews {
		v.FMCreated = v.Created.Format("2006-01-02 15:04:05")
		v.Rurl = v.Url
		v.Url = "/hotnews/detail/" + strconv.Itoa(v.Id) + ".html"
	}

	response["Data"] = hotnews
	response["msg"] = "Success."
	response["code"] = 200
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *HotnewsController) HotDetail() {
	id := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	var hotnews admin.Hotnews
	qs := o.QueryTable(new(admin.Hotnews))
	err := qs.Filter("id", id).RelatedSel().One(&hotnews)
	if err != nil {
		c.Abort("404")
	}

	c.Data["Hotnews"] = hotnews

	var img string
	path := "/static/uploads/" + strconv.Itoa(int(hotnews.Id)) + ".png"
	savepath := "./static/uploads/" + strconv.Itoa(int(hotnews.Id)) + ".png"
	if _, _err := os.Stat(savepath); os.IsNotExist(_err) {
		img = ""
		// img = utils.GetPageImage(savepath,hotnews.Url)
		// if img != "" {
		// 	img = path
		// }
	} else {
		img = path
	}

	if img == "" {
		c.Data["ShowIf"] = true
	} else {
		c.Data["ShowIf"] = false
	}

	c.Data["ImagePath"] = img
	c.TplName = "home/" + c.Template + "/detail.html"

}
