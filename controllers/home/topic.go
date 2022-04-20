package home

import (
	models "go-bbs/models/admin"

	"github.com/astaxie/beego/orm"
)

type TopicController struct {
	BaseController
}

func (c *TopicController) Get() {
	page, _ := c.GetInt("page", 1)
	// 推荐
	o := orm.NewOrm()
	var topic []*models.Topic
	o.QueryTable(new(models.Topic)).Filter("status", 1).OrderBy("-join").Limit(30).Offset((page - 1) * 30).All(&topic)

	type Item struct {
		Title string
		Item  []*models.Topic
	}
	c.Data["AllTopic"] = []Item{
		Item{
			Title: "全部话题",
			Item:  topic,
		},
	}
	c.TplName = "home/" + c.Template + "/square.html"

}
