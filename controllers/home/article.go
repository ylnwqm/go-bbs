package home

import (
	"encoding/json"
	"fmt"
	articleset "go-bbs/models"
	"go-bbs/models/admin"
	"go-bbs/utils"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type ArticleController struct {
	BaseController
}

func (c *ArticleController) News() {

	o := orm.NewOrm()

	var articles []*admin.Article
	qs := o.QueryTable(new(admin.Article))
	qs = qs.Filter("status", 1)
	qs = qs.Filter("cover__isnull", false)
	qs = qs.Filter("Customer__Username__isnull", false)
	qs = qs.Filter("Category__Name__isnull", false)

	var categoryId = 14
	if categoryId != 0 {

		category := new(admin.Category)
		var categorys []*admin.Category
		cqs := o.QueryTable(category)
		cqs = cqs.Filter("status", 1)
		cqs.OrderBy("-sort").All(&categorys)

		ids := utils.CategoryTreeR(categorys, categoryId, 0)

		var cids []int
		cids = append(cids, categoryId)
		for _, v := range ids {
			cids = append(cids, v.Id)
		}

		qs = qs.Filter("Category__ID__in", cids)

	}

	_, err := qs.OrderBy("-id").RelatedSel().Limit(6).Offset(0).All(&articles)
	if err != nil {
		c.Abort("404")
	}

	// for _,v := range articles{
	// 	if v.Cover == ""{
	// 		imgs := utils.FindImg(v.Html)
	// 		if len(imgs) <= 0 {
	// 			continue
	// 		}
	// 		v.Cover = imgs[0]
	// 	}
	// }
	c.Data["index"] = "全站资讯" + " - 霓虹灯下"
	c.Data["Articles"] = articles
	// c.ServeJSON()
	// c.StopRun()
	//fmt.Println(&articles)
	c.TplName = "home/" + c.Template + "/news.html"
}

// 列表
func (c *ArticleController) List() {

	o := orm.NewOrm()
	var setting admin.Setting
	o.QueryTable(new(admin.Setting)).Filter("name", "limit").One(&setting)
	l := setting.Value
	limit, err := strconv.ParseInt(l, 10, 64)
	if err != nil || limit == 0 {
		li, _ := beego.AppConfig.Int64("limit")
		limit = li
	}

	//limit, _ := beego.AppConfig.Int64("limit") // 一页的数量
	page, _ := c.GetInt64("page", 1) // 页数
	offset := (page - 1) * limit     // 偏移量
	// categoryId, _ := c.GetInt("c", 0) //
	cid := c.Ctx.Input.Param(":id")
	categoryId, _ := strconv.Atoi(cid)

	//o := orm.NewOrm()
	article := new(admin.Article)

	var articles []*admin.Article
	qs := o.QueryTable(article)
	qs = qs.Filter("status", 1)
	qs = qs.Filter("Customer__Username__isnull", false)
	qs = qs.Filter("Category__Name__isnull", false)

	if categoryId != 0 {

		category := new(admin.Category)
		var categorys []*admin.Category
		cqs := o.QueryTable(category)
		cqs = cqs.Filter("status", 1)
		cqs.OrderBy("-sort").All(&categorys)

		ids := utils.CategoryTreeR(categorys, categoryId, 0)

		var cids []int
		cids = append(cids, categoryId)
		for _, v := range ids {
			cids = append(cids, v.Id)
		}

		/*c.Data["json"] = cids
		c.ServeJSON()
		c.StopRun()*/

		qs = qs.Filter("Category__ID__in", cids)

	}

	c.Data["CategoryID"] = categoryId
	// 查出当前分类下的所有子分类id

	date := c.GetString("date")
	if date != "" {
		if len(date) == 7 {
			start := ""
			end := ""
			dateNumStr := beego.Substr(date, len("2018-"), 2)
			yearNumStr := beego.Substr(date, len("20"), 2)
			dateNum, _ := strconv.Atoi(dateNumStr)
			yearNum, _ := strconv.Atoi(yearNumStr)

			start = utils.SubString(date, len("2018-01")) + "-01 00:00:00"
			if dateNum >= 12 {
				endYearStr := strconv.Itoa(yearNum + 1)
				end = utils.SubString(date, len("20")) + endYearStr + "-01-01 00:00:00"
			}

			if dateNum < 9 {
				endStr := strconv.Itoa(dateNum + 1)
				end = utils.SubString(date, len("2018-0")) + endStr + "-01 00:00:00"
			}
			if dateNum >= 9 && dateNum < 12 {
				endStr := strconv.Itoa(dateNum + 1)
				end = utils.SubString(date, len("2018-")) + endStr + "-01 00:00:00"
			}

			/*c.Data["json"] = []string{start,end}
			c.ServeJSON()
			c.StopRun()*/

			qs = qs.Filter("created__gte", start)
			qs = qs.Filter("created__lte", end)
			c.Data["Date"] = utils.SubString(start, len("2018-01"))

		} else {
			date = utils.SubString(date, len("2018-01-01"))
			tm, _ := time.Parse("2006-01-02", date)
			unix := tm.Unix() //1566432000

			startFormat := time.Unix(unix, 0).Format("2006-01-02 15:04:05")
			moreUnix, _ := utils.ToInt64(int64(60 * 60 * 24))
			endFormat := time.Unix(unix+moreUnix, 0).Format("2006-01-02 15:04:05")
			start := utils.SubString(startFormat, len("2018-01-01")) + " 00:00:00"
			end := utils.SubString(endFormat, len("2018-01-01")) + " 00:00:00"

			// 刷选
			qs = qs.Filter("created__gte", start)
			qs = qs.Filter("created__lte", end)
			c.Data["Date"] = utils.SubString(start, len("2018-01-01"))
		}
	}

	tag := c.GetString("tag")
	if tag != "" {
		qs = qs.Filter("tag__icontains", tag)
	}

	// 统计
	count, err := qs.Count()
	if err != nil {
		c.Abort("404")
	}

	// 获取数据
	_, err = qs.OrderBy("-id", "-pv").RelatedSel().Limit(limit).Offset(offset).All(&articles)
	if err != nil {
		c.Abort("404")
	}

	//if  c.Template == "app" {
	// for _,v := range articles{
	// 	v.Images = utils.FindImg(v.Html)
	// }
	//}
	c.Data["Data"] = &articles
	c.Data["Paginator"] = utils.GenPaginator(page, limit, count)

	if categoryId == 0 {
		if tag != "" {
			c.Data["index"] = "标签 " + tag + " 的搜索结果 - 霓虹灯下"
		} else {
			c.Data["index"] = "文章大全 - 霓虹灯下社区"
		}

		//c.Data["index"] = "文章大全"
	} else {
		categoryKey := admin.Category{Id: categoryId}
		err = o.Read(&categoryKey)
		c.Data["CategoryName"] = categoryKey.Name
		if err == nil {
			c.Data["index"] = categoryKey.Name + " - 霓虹灯下社区"
		} else {
			c.Data["index"] = "文章大全 - 霓虹灯下社区"
		}
	}

	// 系列文章
	// c.Data["IsArticleSetAll"] = 0
	//fmt.Println(&articles)
	switch c.Template {
	case "nihongdengxia":
		c.TplName = "home/" + c.Template + "/read.html"
	default:
		c.TplName = "home/" + c.Template + "/list.html"
	}
}

// 详情
func (c *ArticleController) Detail() {

	id := c.Ctx.Input.Param(":id")
	viewType := c.GetString("type")

	fnt := c.GetString("from")
	nid, _ := c.GetInt("nid", 0)
	if fnt == "notice" && nid > 0 && c.IsLogin {
		admin.SetReadStatus(&admin.Notice{
			Id: nid,
		}, c.CustomerId)
	}

	// 基础数据
	o := orm.NewOrm()
	article := new(admin.Article)
	var articles []*admin.Article
	qs := o.QueryTable(article)
	err := qs.Filter("id", id).RelatedSel().One(&articles)
	if err != nil {
		c.Abort("404")
	}

	/*c.Data["json"]= &articles
	c.ServeJSON()
	c.StopRun()*/
	tag := articles[0].Tag
	tag = strings.Replace(tag, "，", ",", -1)
	tag_arr := strings.Split(tag, `,`)
	c.Data["Tag"] = tag_arr
	c.Data["Data"] = &articles[0]

	// 系列文章
	var articleSetOne articleset.ArticleSet
	var articleSetTitle string

	o.QueryTable(new(articleset.ArticleSet)).Filter("article_id", articles[0].Id).Filter("status", 1).One(&articleSetOne)
	var articleSetAll []*admin.Article
	if articleSetOne.Id > 0 {
		articleSetTitle = articleSetOne.Title
		var articleSets []*articleset.ArticleSet
		o.QueryTable(new(articleset.ArticleSet)).Filter("title", articleSetOne.Title).Filter("customer_id", articles[0].Customer.Id).Filter("status", 1).All(&articleSets)
		var articlesetIds []int
		articlesetIds = append(articlesetIds, 0)
		for _, v := range articleSets {
			articlesetIds = append(articlesetIds, v.ArticleId)
		}
		//var articleSetAll []*admin.Article
		qssa := o.QueryTable(new(admin.Article))
		qssa = qs.Filter("status", 1)
		qssa = qs.Filter("id__in", articlesetIds)
		_, err = qssa.All(&articleSetAll, "id", "Title", "pv", "like", "review", "Created", "Updated")

	}
	c.Data["ArticleSetAll"] = articleSetAll
	c.Data["ArticleSetTitle"] = articleSetTitle

	// 相关文章
	var relevantArticles []*admin.Article
	rqs := o.QueryTable(new(admin.Article))
	rqs = rqs.Filter("status", 1)
	rqs = rqs.Filter("Customer__Username__isnull", false)
	rqs = rqs.Filter("Category__Name__isnull", false)
	categoryId := articles[0].Category.Id
	if categoryId != 0 {

		category := new(admin.Category)
		var categorys []*admin.Category
		cqs := o.QueryTable(category)
		cqs = cqs.Filter("status", 1)
		cqs.OrderBy("-sort").All(&categorys)

		ids := utils.CategoryTreeR(categorys, categoryId, 0)

		var cids []int
		cids = append(cids, categoryId)
		for _, v := range ids {
			cids = append(cids, v.Id)
		}

		/*c.Data["json"] = cids
		c.ServeJSON()
		c.StopRun()*/

		rqs = rqs.Filter("Category__ID__in", cids)

	}
	// 获取数据
	_, err = rqs.OrderBy("-recommend", "-id", "-pv").RelatedSel().Limit(5).All(&relevantArticles)
	if err != nil {
		c.Abort("404")
	}

	c.Data["RelevantArticles"] = &relevantArticles

	if beego.AppConfig.String("view") == "default" {
		var listData = make(map[string][]*admin.Article)
		var list []*admin.Article
		_, err = o.QueryTable(article).Filter("status", 1).Filter("User__Name__isnull", false).Filter("Category__Name__isnull", false).OrderBy("id").RelatedSel().All(&list, "id", "title")

		for _, v := range list {
			listData[v.Category.Name] = append(listData[v.Category.Name], v)
		}
		c.Data["List"] = &listData
		articleId, _ := strconv.Atoi(id)
		c.Data["ArticleId"] = articleId
		/*c.Data["json"]= &listData
		c.ServeJSON()
		c.StopRun()*/

	}

	var other admin.Other

	if &articles[0].Other != nil {
		json.Unmarshal([]byte(articles[0].Other), &other)
	}

	other.SubjectInfo = strings.Replace(other.SubjectInfo, "\n", "<br>", -1)
	c.Data["index"] = &articles[0].Title
	c.Data["Other"] = other

	if viewType == "single" {
		c.TplName = "home/" + c.Template + "/doc.html"
	} else {
		c.TplName = "home/" + c.Template + "/detail.html"
	}
	//c.TplName = "home/nihongdengxia/review.html"
}

// 统计访问量
func (c *ArticleController) Pv() {

	ids := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(ids)
	/*c.Data["json"] = c.Input()
	c.ServeJSON()
	c.StopRun()*/

	response := make(map[string]interface{})

	o := orm.NewOrm()

	article := admin.Article{Id: id}
	if o.Read(&article) == nil {
		article.Pv = article.Pv + 1

		valid := validation.Validation{}
		valid.Required(article.Id, "Id")
		if valid.HasErrors() {
			// 如果有错误信息，证明验证没通过
			// 打印错误信息
			for _, err := range valid.Errors {
				//log.Println(err.Key, err.Message)
				response["msg"] = "Error."
				response["code"] = 500
				response["err"] = err.Key + " " + err.Message
				c.Data["json"] = response
				c.ServeJSON()
				c.StopRun()
			}
		}

		if _, err := o.Update(&article); err == nil {
			response["msg"] = "Success."
			response["code"] = 200
			response["id"] = id
		} else {
			response["msg"] = "Error."
			response["code"] = 500
			response["err"] = err.Error()
		}
	} else {
		response["msg"] = "Error."
		response["code"] = 500
		response["err"] = "ID 不能为空！"
	}

	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

// 评论
func (c *ArticleController) Review() {

	response := make(map[string]interface{})
	customer := c.GetSession("Customer")
	var User admin.Customer
	if customer != nil {
		User = customer.(admin.Customer)
	} else {
		response["msg"] = "请先登录！"
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	aid, _ := c.GetInt("aid")
	review := c.GetString("review")
	is_push_bbs := c.GetString("is_push_bbs")

	o := orm.NewOrm()
	reviewsMd := admin.Review{
		Name:      User.Username,
		Review:    template.HTMLEscapeString(review),
		Site:      "",
		ArticleId: aid,
		Status:    1,
		Customer:  &admin.Customer{Id: User.Id},
	}

	valid := validation.Validation{}
	//valid.Required(reviewsMd.Name, "Name")
	valid.Required(reviewsMd.Review, "Review")
	valid.Required(reviewsMd.ArticleId, "ArticleId")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			//log.Println(err.Key, err.Message)
			response["msg"] = "新增失败！" + err.Key + " " + err.Message
			response["code"] = 500
			response["err"] = err.Key + " " + err.Message
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		}
	}

	// 更新评论数量
	o.Begin()
	article := admin.Article{Id: aid}
	o.Read(&article)
	article.Review = article.Review + 1
	o.Update(&article)

	if is_push_bbs == "1" {
		admin.AddBbs(&admin.Bbs{
			Customer:  &admin.Customer{Id: c.CustomerId},
			Content:   reviewsMd.Review,
			ArticleId: aid,
			Created:   time.Now(),
			Updated:   time.Now(),
			Status:    1,
		})
	}

	if id, err := o.Insert(&reviewsMd); err == nil {
		response["msg"] = "新增成功！"
		response["code"] = 200
		response["id"] = id
		o.Commit()
	} else {
		response["msg"] = "新增失败！"
		response["code"] = 500
		response["err"] = err.Error()
	}

	type EmailData struct {
		Review  admin.Review
		Article admin.Article
	}

	// 发送邮件
	customerEmail, err := admin.GetCustomerById(int(article.Customer.Id))
	if customerEmail.Email != "" || err != nil {
		flag := utils.SendEmail(utils.Email{
			From: "1920853199@qq.com",
			To:   customerEmail.Email,
			Header: map[string]string{
				"Subject": User.Username + "评论了您的文章",
			},
			Template: "views/home/nihongdengxia/comment-email.html",
			Data: EmailData{
				Review:  reviewsMd,
				Article: article,
			},
		})
		if flag == nil {
			utils.SendNotic(utils.NoticeMessage{
				Title:     "评论了你的文章",
				SendId:    User.Id,
				ReceiveId: article.Customer.Id,
				Username:  User.Username,
				UserUrl:   fmt.Sprintf(`/user/%d`, User.Id),
				Content:   article.Title,
				Replay:    reviewsMd.Review,
				Url:       fmt.Sprintf(`/detail/%d.html`, article.Id),
				Date:      time.Now(),
				Cover:     User.Image,
			})

			fmt.Println("评论邮件发送成功！")
		} else {
			fmt.Println("评论邮件发送失败！")
		}
	}
	// Id        int    `orm:"column(id);auto" description:"ID"`
	// SendId    int    `orm:"column(send_id)" description:"发送人"`
	// ReceiveId int    `orm:"column(receive_id)" description:"接收人"`
	// Type      int    `orm:"column(type)" description:"消息类型"`
	// Title     string `orm:"column(title);size(50)" description:"标题"`
	// Content   string `orm:"column(content);size(50)" description:"消息内容"`
	// Status    int    `orm:"column(status);null" description:"1未读，2已读禁用"`

	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *ArticleController) ReviewList() {

	id := c.Ctx.Input.Param(":id")
	//limit, _ := beego.AppConfig.Int64("limit") // 一页的数量
	limit := int64(20)
	page, _ := c.GetInt64("page", 1) // 页数
	offset := (page - 1) * limit     // 偏移量
	response := make(map[string]interface{})

	o := orm.NewOrm()
	review := new(admin.Review)

	var reviews []*admin.Review
	qs := o.QueryTable(review)
	qs = qs.Filter("status", 1)
	qs = qs.Filter("article_id", id)

	// 获取数据
	_, err := qs.OrderBy("-id").RelatedSel().Limit(limit).Offset(offset).All(&reviews)

	if err != nil {
		response["msg"] = "Error."
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	// 统计
	count, err := qs.Count()
	if err != nil {
		response["msg"] = "Error."
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	response["Data"] = &reviews
	response["Paginator"] = utils.GenPaginator(page, limit, count)

	response["msg"] = "Success."
	response["code"] = 200
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()

}

func (c *ArticleController) Like() {

	response := make(map[string]interface{})
	ip := c.Ctx.Input.IP()
	id, _ := c.GetInt("id")

	o := orm.NewOrm()
	qs := o.QueryTable(new(admin.Log))
	// fmt.Println(ip)

	qs = qs.Filter("ip", ip)
	qs = qs.Filter("create__gte", beego.Date(time.Now(), "Y-m-d 00:00:00"))
	qs = qs.Filter("create__lte", beego.Date(time.Now(), "Y-m-d H:i:s"))
	qs = qs.Filter("page", "/article/like"+strconv.Itoa(id))

	count, e := qs.Count()

	if e != nil {
		response["msg"] = "Error."
		response["code"] = 500
		response["err"] = e.Error()
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}
	if count > 1 {
		response["msg"] = "Error."
		response["code"] = 500
		response["err"] = "亲，点赞过了，明天再来哦！"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	article := admin.Article{Id: id}
	if o.Read(&article) == nil {
		article.Like = article.Like + 1

		valid := validation.Validation{}
		valid.Required(article.Id, "Id")
		if valid.HasErrors() {
			// 如果有错误信息，证明验证没通过
			// 打印错误信息
			for _, err := range valid.Errors {
				//log.Println(err.Key, err.Message)
				response["msg"] = "Error."
				response["code"] = 500
				response["err"] = err.Key + " " + err.Message
				c.Data["json"] = response
				c.ServeJSON()
				c.StopRun()
			}
		}

		if _, err := o.Update(&article); err == nil {
			response["msg"] = "Success."
			response["code"] = 200
			response["like"] = article.Like
		} else {
			response["msg"] = "Error."
			response["code"] = 500
			response["err"] = err.Error()
		}
	} else {
		response["msg"] = "Error."
		response["code"] = 500
		response["err"] = "ID 不能为空！"
	}

	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}
