package home

import (
	"go-bbs/models"
	"go-bbs/models/admin"

	"github.com/astaxie/beego/orm"
)

type ArticleSetController struct {
	BaseController
}

func (c *ArticleSetController) GetArticleBySet() {

	response := make(map[string]interface{})
	set := c.GetString("set")
	cid, _ := c.GetInt("cid")
	if set == "" || cid == 0 {
		response["msg"] = "非法操作"
		response["code"] = 200
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	var articleSets []*models.ArticleSet
	o := orm.NewOrm()
	o.QueryTable(new(models.ArticleSet)).Filter("title", set).Filter("customer_id", cid).Filter("status", 1).All(&articleSets)

	var articleIds []int
	articleIds = append(articleIds, 0)
	for _, v := range articleSets {
		articleIds = append(articleIds, v.ArticleId)
	}

	// fmt.Println(articleIds)
	var articles []*admin.Article
	qs := o.QueryTable(new(admin.Article))
	qs = qs.Filter("status", 1)
	qs = qs.Filter("id__in", articleIds)
	_, err := qs.All(&articles, "id", "Title", "pv", "like", "review", "Created", "Updated")
	if err != nil {
		c.Abort("404")
	}

	response["Data"] = articles
	response["msg"] = "Success."
	response["code"] = 200
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}
