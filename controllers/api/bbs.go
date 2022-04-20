package api

import (
	models "go-bbs/models/admin"
	"go-bbs/utils"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type BbsController struct {
	BaseController
}

func (ctl *BbsController) Bbs() {
	// if !ctl.IsLogin {
	// 	ctl.Redirect("/login.html", 302)
	// }

	limit := int64(50)
	page, _ := ctl.GetInt64("page", 1) // 页数
	offset := (page - 1) * limit       // 偏移量
	id := ctl.Ctx.Input.Param(":id")
	category_id, _ := strconv.Atoi(id)

	o := orm.NewOrm()
	var bbs []*models.Bbs
	qs := o.QueryTable(new(models.Bbs))
	qs = qs.Filter("status", 1)
	if category_id > 0 {
		var t []models.Topic
		o.QueryTable(new(models.Topic)).Filter("status", 1).Filter("category_id", category_id).All(&t)
		var topic_id []int
		for _, v := range t {
			topic_id = append(topic_id, v.Id)
		}
		if len(topic_id) > 0 {
			qs = qs.Filter("topic_id__in", topic_id)
		} else {
			qs = qs.Filter("topic_id__in", -1)
		}
	}
	qs = qs.Filter("Customer__Username__isnull", false)

	count, err := qs.Count()
	if err != nil {
		ctl.Abort("404")
	}

	// 获取数据
	_, err = qs.OrderBy("-id").RelatedSel().Limit(limit).Offset(offset).All(&bbs)
	if err != nil {
		panic(err)
	}

	for _, v := range bbs {

		//<a href="#">(展开)</a>
		v.Content = utils.Subbbs(v.Content, 200, `/bbs/detail/`+strconv.Itoa(v.Id)+`.html`)
		orm.NewOrm().LoadRelated(v, "Images")
		//orm.NewOrm().LoadRelated(v, "BbsReview")
		var review []*models.BbsReview
		o.QueryTable(new(models.BbsReview)).RelatedSel().Filter("Status", 1).Filter("bbs_id", v.Id).All(&review)
		v.BbsReview = review
		// v.IsLike = o.QueryTable(new(models.BbsLike)).Filter("bbs_id", v.Id).Filter("customer_id", ctl.CustomerId).Exist()
		v.IsLike = o.QueryTable(new(models.BbsLike)).Filter("bbs_id", v.Id).Filter("customer_id", 1).Exist()

		if v.TopicId != 0 {
			var topic models.Topic
			err = o.QueryTable(new(models.Topic)).Filter("id", v.TopicId).Filter("status", 1).One(&topic)
			v.Topic = &topic
		}

		//o.QueryTable(new(models.BbsReview)).RelatedSel().Filter("Status", 1).Filter("bbs_id", v.Id).All(&review)

	}
	/*
		// 评论不要暂时
		var review = make(map[int][]utils.ReviewTree)
		for _, v := range bbs {
			review[v.Id] = utils.ReviewTreeR(v.BbsReview, 0, 0, &models.BbsReview{})
		}

		ctl.Data["Review"] = &review
	*/
	ctl.Data["Data"] = &bbs
	ctl.Data["Paginator"] = utils.GenPaginator(page, limit, count)

	/*
		// 主题不要暂时
		var topic []*admin.Topic
		o.QueryTable(new(admin.Topic)).Filter("status", 1).OrderBy("-join").Limit(10).All(&topic)
		ctl.Data["Topic"] = topic

		// 活跃用户不要暂时
		var hotUser []*models.Customer
		o.QueryTable(new(models.Customer)).Filter("Email__isnull", false).OrderBy("-id").Limit(5).Offset(0).All(&hotUser)
		for _, v := range hotUser {
			v.IsFans = models.IsFans(ctl.CustomerId, v.Id)
		}

		ctl.Data["hotUser"] = hotUser
	*/

	/*
		// 动态数量
		bbsCount, _ := o.QueryTable(new(models.Bbs)).Filter("status", 1).Filter("customer_id", ctl.CustomerId).Count()
		ctl.Data["BbsCount"] = bbsCount
		// 访问量
		ctl.Data["Visit"] = models.GetVisitByCustomerId(ctl.CustomerId)
	*/
	ctl.Data["json"] = bbs
	ctl.ServeJSON()
	ctl.StopRun()

}

// func (ctl *BbsController) CRouters() {
// 	utils.CreateRouters()
// }
