package admin

import (
	"go-bbs/models/admin"
	"go-bbs/utils"

	"github.com/astaxie/beego/orm"
)

type LogController struct {
	BaseController
}

func (c *LogController) GetAll() {
	o := orm.NewOrm()

	limit := int64(15)
	page, _ := c.GetInt64("page", 1) // 页数
	offset := (page - 1) * limit     // 偏移量

	start := c.GetString("start")
	end := c.GetString("end")

	referer := c.GetString("referer")

	c.Data["Start"] = start
	c.Data["End"] = end
	c.Data["Referer"] = referer
	//o := orm.NewOrm()
	var logs []*admin.Log
	qs := o.QueryTable(new(admin.Log))

	// 开始时间
	if start != "" {
		qs = qs.Filter("create__gte", start)
	}

	// 结束时间
	if end != "" {
		qs = qs.Filter("create__lte", end)
	}
	// 标题
	if referer != "" {
		qs = qs.Filter("referer__icontains", referer)
	}
	// 获取数据
	qs.OrderBy("-id").Limit(limit).Offset(offset).All(&logs)

	// 统计
	count, _ := qs.Count()

	c.Data["Data"] = &logs
	c.Data["Paginator"] = utils.GenPaginator(page, limit, count)

	c.TplName = "admin/logs.html"
}
