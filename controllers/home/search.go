package home

import (
	"fmt"
	"go-bbs/models/admin"
	models "go-bbs/models/admin"
	"go-bbs/utils"

	"github.com/astaxie/beego/orm"
)

type SearchController struct {
	BaseController
}
type SearchData struct {
	Title       string
	Url         string
	Describe    string
	Time        string
	Category    string
	CategoryUrl string
	Pv          int
}

func (c *SearchController) Search() {
	q := c.GetString("q")
	c.Data["q"] = q
	ctype, _ := c.GetInt("c")
	page, _ := c.GetInt64("page", 1)
	limit := int64(10)
	offset := (page - 1) * limit

	o := orm.NewOrm()

	var sd []SearchData
	var count int64
	// 文章
	if ctype == 0 || ctype == 1 {
		var articles []*admin.Article
		var categorys []*admin.Category
		cqs := o.QueryTable(new(admin.Category))
		cqs = cqs.Filter("name__icontains", q)
		cqs.OrderBy("-sort").All(&categorys)
		var cids []int
		cids = append(cids, 0)
		for _, v := range categorys {
			cids = append(cids, v.Id)
		}
		cond := orm.NewCondition()
		cond = cond.AndCond(cond.And("status", 1).And("Customer__Username__isnull", false).And("Category__Name__isnull", false).And("title__icontains", q)).OrCond(cond.And("status", 1).And("Category__ID__in", cids)).OrCond(cond.And("status", 1).And("Customer__Username__isnull", false).And("Category__Name__isnull", false).And("tag__icontains", q))
		qs := o.QueryTable(new(admin.Article))
		qs = qs.SetCond(cond)
		_, err := qs.OrderBy("-recommend", "-pv", "-id").RelatedSel().Limit(limit).Offset(offset).All(&articles)

		if err != nil {
			c.Data["error"] = err.Error()
			return
		}

		// 统计
		acount, err := qs.Count()
		if err != nil {
			c.Abort("404")
		}
		count += acount
		for _, v := range articles {
			sd = append(sd, SearchData{
				Title:       v.Title,
				Url:         fmt.Sprintf("/detail/%d.html", v.Id),
				Describe:    v.Remark,
				Time:        v.Created.Format("2006-01-02 15:04:05"),
				Category:    v.Category.Name,
				CategoryUrl: fmt.Sprintf("/article/category/%d.html", v.Category.Id),
				Pv:          v.Pv,
			})
		}
	}
	// 用户
	if ctype == 0 || ctype == 5 {
		var sUser []*models.Customer
		susermodel := o.QueryTable(new(models.Customer)).Filter("Email__isnull", false).Filter("username__icontains", q).Limit(limit).Offset(offset)
		susermodel.All(&sUser)
		ucount, err := susermodel.Count()
		if err != nil {
			c.Abort("404")
		}
		count += ucount
		if len(sUser) > 0 {
			for _, v := range sUser {
				sd = append(sd, SearchData{
					Title:    v.Username,
					Url:      fmt.Sprintf("/user/%d", v.Id),
					Describe: v.Signature,
					Time:     v.Created.Format("2006-01-02 15:04:05"),
				})
			}
		}
	}

	// 话题
	if ctype == 0 || ctype == 3 {
		var topic []*models.Topic
		tmodel := o.QueryTable(new(models.Topic)).Filter("status", 1).Filter("content__icontains", q).Limit(limit).Offset(offset)
		tmodel.All(&topic)
		tcount, err := tmodel.Count()
		if err != nil {
			c.Abort("404")
		}
		count += tcount
		if len(topic) > 0 {
			for _, v := range topic {
				sd = append(sd, SearchData{
					Title:    v.Content,
					Url:      fmt.Sprintf("/topic/detail/%d.html", v.Id),
					Describe: fmt.Sprintf("有 %d 人参与了讨论", v.Join),
					Time:     v.Created.Format("2006-01-02 15:04:05"),
				})
			}
		}
	}

	// 帖子
	if ctype == 0 || ctype == 4 {
		var bbs []*models.Bbs
		bbsqs := o.QueryTable(new(models.Bbs))
		bbsqs = bbsqs.Filter("status", 1)
		bbsqs = bbsqs.Filter("Customer__Username__isnull", false)
		// 获取数据
		bbsqs = bbsqs.Filter("content__icontains", q).Limit(limit).Offset(offset)
		bbsqs.All(&bbs)
		bbscount, err := bbsqs.Count()
		if err != nil {
			c.Abort("404")
		}
		count += bbscount
		if len(bbs) > 0 {
			for _, v := range bbs {
				sd = append(sd, SearchData{
					Title:    v.Content,
					Url:      fmt.Sprintf("/bbs/detail/%d.html", v.Id),
					Describe: fmt.Sprintf("有 %d 人觉得很赞，有 %d 参与了评论", v.Like, v.Review),
					Time:     v.Created.Format("2006-01-02 15:04:05"),
				})
			}
		}
	}

	c.Data["Data"] = sd
	c.Data["Paginator"] = utils.GenPaginator(page, limit, count)

	c.Data["count"] = count
	c.TplName = "home/" + c.Template + "/search.html"
}
