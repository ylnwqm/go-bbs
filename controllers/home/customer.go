package home

import (
	"fmt"
	articleset "go-bbs/models"
	"go-bbs/models/admin"
	models "go-bbs/models/admin"
	"go-bbs/utils"
	"path"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

// CustomerController operations for Customer
type CustomerController struct {
	BaseController
}

func (c *CustomerController) List() {

	limit := int64(30)
	page, _ := c.GetInt64("page", 1) // 页数
	offset := (page - 1) * limit     // 偏移量

	var customers []*models.Customer
	o := orm.NewOrm()
	o.QueryTable(new(models.Customer)).OrderBy("-id").Limit(limit).Offset(offset).All(&customers)

	count, _ := o.QueryTable(new(models.Customer)).Count()
	c.Data["Members"] = &customers
	c.Data["Paginator"] = utils.GenPaginator(page, limit, count)
	c.TplName = "home/" + c.Template + "/members.html"
}

func (c *CustomerController) UserInfo() {

	// if !c.IsLogin {
	// 	 c.Redirect("/login.html", 302)
	// }
	fnt := c.GetString("from")
	nid, _ := c.GetInt("nid", 0)
	if fnt == "notice" && nid > 0 && c.IsLogin {
		admin.SetReadStatus(&admin.Notice{
			Id: nid,
		}, c.CustomerId)
	}

	id := c.Ctx.Input.Param(":id")
	userId, _ := strconv.Atoi(id)

	o := orm.NewOrm()
	var homeUser models.Customer
	err := o.QueryTable(new(models.Customer)).Filter("id", userId).One(&homeUser)
	if err != nil {
		c.Abort("404")
	}
	c.Data["HomeUser"] = homeUser

	limit := int64(50)
	page, _ := c.GetInt64("page", 1) // 页数
	offset := (page - 1) * limit     // 偏移量

	var bbs []*models.Bbs
	qs := o.QueryTable(new(models.Bbs))
	qs = qs.Filter("status", 1)
	qs = qs.Filter("Customer__Username__isnull", false)
	qs = qs.Filter("customer_id", userId)
	count, err := qs.Count()
	if err != nil {
		c.Abort("404")
	}

	bbsCount, _ := qs.Count()
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
		v.IsLike = o.QueryTable(new(models.BbsLike)).Filter("bbs_id", v.Id).Filter("customer_id", c.CustomerId).Exist()

		if v.TopicId != 0 {
			var topic models.Topic
			err = o.QueryTable(new(models.Topic)).Filter("id", v.TopicId).Filter("status", 1).One(&topic)
			v.Topic = &topic
		}

		if v.ArticleId != 0 {
			var article models.Article
			err = o.QueryTable(new(models.Article)).Filter("id", v.ArticleId).Filter("status", 1).One(&article)
			v.Article = &article
		}

		if len(v.Images) > 0 {
			for _, vv := range v.Images {
				ext := path.Ext(vv.Url)
				if ext == ".mp4" || ext == ".avi" || ext == ".mov" {
					vv.IsVideo = true
				}
			}
		}
		o.QueryTable(new(models.BbsReview)).RelatedSel().Filter("Status", 1).Filter("bbs_id", v.Id).All(&review)

	}
	var review = make(map[int][]utils.ReviewTree)
	for _, v := range bbs {
		review[v.Id] = utils.ReviewTreeR(v.BbsReview, 0, 0, &models.BbsReview{})
	}

	isFans := false
	if c.IsLogin {
		isFans = models.IsFans(c.CustomerId, userId)
	}

	// 文章
	var articles []*admin.Article
	aqs := o.QueryTable(new(admin.Article))
	aqs = aqs.Filter("Customer__Username__isnull", false)
	aqs = aqs.Filter("Category__Name__isnull", false)
	aqs = aqs.Filter("customer_id", userId)
	// 获取数据
	aqs.OrderBy("-id").Limit(limit).Offset(offset).All(&articles)
	acount, _ := aqs.Count()
	c.Data["Acount"] = acount
	c.Data["Articles"] = articles

	c.Data["Fans"] = models.GetFans(homeUser.Id)
	c.Data["Focus"] = admin.GetFocus(homeUser.Id)

	c.Data["bbsCount"] = bbsCount
	c.Data["IsFans"] = isFans
	c.Data["Review"] = &review
	c.Data["Data"] = &bbs
	c.Data["Paginator"] = utils.GenPaginator(page, limit, count)
	c.Data["IsTopic"] = false
	c.Data["ReviewShow"] = false

	// 文集
	articleSet, _ := articleset.GetArticleSetCustomerId(homeUser.Id)
	c.Data["ArticleSet"] = articleSet

	for _, v := range articleSet {
		qs := o.QueryTable(new(articleset.ArticleSet))
		qs = qs.Filter("customer_id", homeUser.Id)
		qs = qs.Filter("status", 1)
		qs = qs.Filter("title", v.Title)
		v.Count, _ = qs.Count()
		//qs = qs.Filter("title", a.Id)
	}

	c.TplName = "home/" + c.Template + "/user_info.html"
}

func (c *CustomerController) Profile() {

	if !c.IsLogin {
		c.Redirect("/login.html"+"?referrer="+c.Ctx.Input.URL(), 302)
		//c.Redirect("/login.html", 302)
	}
	customer, _ := admin.GetCustomerById(c.CustomerId)
	c.Data["CustomerData"] = customer
	c.Data["CustomerPage"] = "Profile"
	c.TplName = "home/" + c.Template + "/profile.html"
}

func (c *CustomerController) Put() {

	response := make(map[string]interface{})
	if !c.IsLogin {
		response["code"] = 500
		response["msg"] = "请先登录"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}
	v := admin.Customer{
		Id:       c.CustomerId,
		Nickname: c.GetString("nickname"),
		Username: c.GetString("username"),
		//Email:     c.GetString("email"),
		Phone:     c.GetString("phone"),
		Image:     c.GetString("image"),
		Signature: c.GetString("signature"),
		Url:       c.GetString("url"),
	}

	if err := admin.UpdateCustomerById(&v); err == nil {

		o := orm.NewOrm()
		customerSeesion := admin.Customer{Id: c.CustomerId}
		err := o.Read(&customerSeesion)
		if err != nil {
			response["code"] = 500
			response["msg"] = err.Error()
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		}
		c.SetSession("Customer", customerSeesion)
		response["code"] = 200
		response["msg"] = "修改成功！"
	} else {
		response["code"] = 500
		response["msg"] = "用户名或者密码不能为空！"
	}
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *CustomerController) Fans() {
	response := make(map[string]interface{})
	if !c.IsLogin {
		response["code"] = 500
		response["msg"] = "请先登录"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}
	focus_id, _ := c.GetInt("focus_id")
	customer_id := c.CustomerId
	if focus_id < 1 {
		response["code"] = 500
		response["msg"] = "非法focus_id"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	status, flag := models.AddFocus(customer_id, focus_id)
	if flag && status == 1 {

		utils.SendNotic(utils.NoticeMessage{
			Title:     "关注了你",
			SendId:    c.CustomerId,
			ReceiveId: focus_id,
			Username:  c.Username,
			UserUrl:   fmt.Sprintf(`/user/%d`, c.CustomerId),
			Content:   "",
			Replay:    "",
			Url:       fmt.Sprintf(`/user/%d`, focus_id),
			Date:      time.Now(),
			Cover:     c.Cover,
		})

		response["status"] = status
		response["code"] = 200
		response["msg"] = "操作成功！"
	} else {
		response["code"] = 200
		response["msg"] = "操作失败！"
	}
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *CustomerController) Visit() {
	response := make(map[string]interface{})

	focus_id, _ := c.GetInt("focus_id")
	customer_id := c.CustomerId
	if focus_id < 1 {
		response["code"] = 500
		response["msg"] = "非法focus_id"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	status, flag := models.AddFocus(customer_id, focus_id)
	if flag {
		response["status"] = status
		response["code"] = 200
		response["msg"] = "操作成功！"
	} else {
		response["code"] = 200
		response["msg"] = "操作失败！"
	}
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}
