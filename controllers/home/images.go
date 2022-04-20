package home

import (
	models "go-bbs/models/admin"

	"github.com/astaxie/beego/orm"
)

type ImagesController struct {
	BaseController
}

func (c *ImagesController) List() {
	o := orm.NewOrm()
	var images []*models.BbsImages
	o.QueryTable(new(models.BbsImages)).OrderBy("-id").RelatedSel().Limit(40).All(&images)
	c.Data["Images"] = images
	c.TplName = "home/" + c.Template + "/images.html"
}
