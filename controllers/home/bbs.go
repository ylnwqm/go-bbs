package home

import (
	"fmt"
	"go-bbs/models/admin"
	models "go-bbs/models/admin"
	"go-bbs/utils"
	histoday "go-bbs/utils/hotnews"
	"html/template"
	"path"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/astaxie/beego/orm"
)

type BbsController struct {
	BaseController
}

func (ctl *BbsController) BbsCreate() {
	if !ctl.IsLogin {
		ctl.Redirect("/login.html"+"?referrer="+ctl.Ctx.Input.URL(), 302)
		//ctl.Redirect("/login.html", 302)
	}
	o := orm.NewOrm()
	var topicContent []*models.Topic
	o.QueryTable(new(models.Topic)).Filter("status", 1).All(&topicContent)
	ctl.Data["CTopic"] = &topicContent
	//fmt.Printf("%v\n",topicContent)
	ctl.TplName = "home/" + ctl.Template + "/bbs_create.html"
}

func (ctl *BbsController) Bbs() {
	if !ctl.IsLogin {
		//ctl.Redirect("/login.html", 302)
	}
	var focus []models.Customer
	t := ctl.GetString("t")
	if t == "" || t == "all" {
		ctl.Data["BbsT"] = "all"
	} else if t == "follow" {
		if !ctl.IsLogin {
			ctl.Redirect("/bbs.html", 302)
		}
		focus = models.GetFocus(ctl.CustomerId)
		ctl.Data["BbsT"] = "follow"
	} else {
		ctl.Data["BbsT"] = "all"
	}
	o := orm.NewOrm()

	if ctl.Template == "baidu" {
		((*HotnewsController)(unsafe.Pointer(ctl))).List()

		var hotnews []*admin.Hotnews
		o.QueryTable(new(admin.Hotnews)).RelatedSel().OrderBy("-id").Limit(100).All(&hotnews)

		for _, v := range hotnews {
			//fmt.Println(v.Category)
			v.Rurl = v.Url
			v.Url = "/dailyhot/detail/" + strconv.Itoa(v.Id) + ".html"
		}

		ctl.Data["Hotnews"] = hotnews

		h := histoday.Histoday{}
		Histoday := h.Get("https://hao.360.com/histoday/")
		firstHistoday := Histoday[0:6]
		moreHistoday := Histoday[6:]
		ctl.Data["Histoday"] = firstHistoday
		ctl.Data["MoreHistoday"] = moreHistoday
		p := histoday.Hupu{}
		ctl.Data["Match"] = p.GetMatch("https://www.hupu.com/")
	}

	limit := int64(20)
	page, _ := ctl.GetInt64("page", 1) // 页数
	offset := (page - 1) * limit       // 偏移量
	id := ctl.Ctx.Input.Param(":id")
	category_id, _ := strconv.Atoi(id)

	var bbs []*models.Bbs
	qs := o.QueryTable(new(models.Bbs))
	qs = qs.Filter("status", 1)
	var description string
	var CategorydName string
	if category_id > 0 {

		qs = qs.Filter("category_id__in", category_id)

		var cate models.Category

		cate_err := o.QueryTable(new(models.Category)).Filter("id", category_id).One(&cate)
		//fmt.SPrintf("%d",category_id)
		ctl.Data["CategorydImage"] = "/static/images/category/" + fmt.Sprintf("%d", category_id) + ".jpeg"
		if cate_err == nil {
			ctl.Data["SeoTitle"] = "有关" + cate.Name + "的帖子"
		} else {
			ctl.Data["SeoTitle"] = "大厅 - 霓虹灯下社区"
		}
		CategorydName = cate.Name
		description = cate.Description
	} else {
		ctl.Data["SeoTitle"] = "大厅 - 霓虹灯下社区"
	}
	ctl.Data["CategorydName"] = CategorydName
	ctl.Data["CategorydDescription"] = description
	qs = qs.Filter("Customer__Username__isnull", false)
	if t == "follow" && len(focus) > 0 {
		var focus_ids []int
		for _, v := range focus {
			focus_ids = append(focus_ids, v.Id)
		}
		qs = qs.Filter("customer_id__in", focus_ids)
	}

	count, err := qs.Count()
	if err != nil {
		ctl.Abort("404")
	}

	// 获取数据
	_, err = qs.OrderBy("-updated").RelatedSel().Limit(limit).Offset(offset).All(&bbs)
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
		if ctl.IsLogin {
			v.IsLike = o.QueryTable(new(models.BbsLike)).Filter("bbs_id", v.Id).Filter("customer_id", ctl.CustomerId).Exist()
		} else {
			v.IsLike = false
		}

		if v.TopicId != 0 {
			var topic models.Topic
			err = o.QueryTable(new(models.Topic)).Filter("id", v.TopicId).Filter("status", 1).One(&topic)
			v.Topic = &topic
		}

		if v.ArticleId != 0 {
			var article models.Article
			err = o.QueryTable(new(models.Article)).Filter("id", v.ArticleId).Filter("status", 1).One(&article)
			if article.Cover == "" {
				img := utils.FindImg(article.Html)
				if len(img) > 0 {
					article.Cover = img[0]
				} else {
					article.Cover = "/static/images/bitimage.png"
				}

			}
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
		//o.QueryTable(new(models.BbsReview)).RelatedSel().Filter("Status", 1).Filter("bbs_id", v.Id).All(&review)

	}
	var review = make(map[int][]utils.ReviewTree)
	for _, v := range bbs {
		review[v.Id] = utils.ReviewTreeR(v.BbsReview, 0, 0, &models.BbsReview{})
	}

	ctl.Data["CategoryId"] = category_id
	ctl.Data["Review"] = &review
	ctl.Data["Data"] = &bbs
	ctl.Data["Paginator"] = utils.GenPaginator(page, limit, count)

	// 主题
	var topic []*admin.Topic
	o.QueryTable(new(admin.Topic)).Filter("status", 1).OrderBy("-join").Limit(10).All(&topic)
	ctl.Data["Topic"] = topic

	// 活跃用户
	var hotUser []*models.Customer
	o.QueryTable(new(models.Customer)).Filter("Email__isnull", false).OrderBy("-integral").Limit(10).Offset(0).All(&hotUser)
	for _, v := range hotUser {
		if ctl.IsLogin {
			v.IsFans = models.IsFans(ctl.CustomerId, v.Id)
		} else {
			v.IsFans = false
		}
	}

	ctl.Data["hotUser"] = hotUser
	if ctl.IsLogin {
		// 动态数量
		bbsCount, _ := o.QueryTable(new(models.Bbs)).Filter("status", 1).Filter("customer_id", ctl.CustomerId).Count()
		ctl.Data["BbsCount"] = bbsCount
		// 访问量
		ctl.Data["Visit"] = models.GetVisitByCustomerId(ctl.CustomerId)
	}
	// ctl.Data["json"] = bbs
	// ctl.ServeJSON()
	// ctl.StopRun()
	// 评论框是否显示
	ctl.Data["ReviewShow"] = false
	ctl.Data["IsTopic"] = false
	bbsview := ctl.Ctx.GetCookie("bbsview")

	if bbsview == "program" {
		if ctl.Template == "app" {
			ctl.Data["BBSVIEW"] = "program"
			ctl.TplName = "home/" + ctl.Template + "/bbs.html"
		} else {
			ctl.Data["BBSVIEW"] = "bbs"
			ctl.TplName = "home/" + ctl.Template + "/list_bbs.html"
		}
	} else {
		ctl.Data["BBSVIEW"] = "program"
		ctl.TplName = "home/" + ctl.Template + "/bbs.html"
	}
}

func (ctl *BbsController) BbsDetail() {
	if !ctl.IsLogin {
		//ctl.Redirect("/login.html", 302)
	}

	fnt := ctl.GetString("from")
	nid, _ := ctl.GetInt("nid", 0)
	if fnt == "notice" && nid > 0 && ctl.IsLogin {
		admin.SetReadStatus(&admin.Notice{
			Id: nid,
		}, ctl.CustomerId)
	}

	id := ctl.Ctx.Input.Param(":id")
	bbs_id, _ := strconv.Atoi(id)

	o := orm.NewOrm()
	var bbss []*models.Bbs
	qs := o.QueryTable(new(models.Bbs))
	qs = qs.Filter("status", 1)
	qs = qs.Filter("Customer__Username__isnull", false)
	qs = qs.Filter("id", bbs_id)

	// 获取数据
	_, err := qs.OrderBy("-id").RelatedSel().All(&bbss)
	if err != nil {
		panic(err)
	}

	bbs := bbss[0]

	orm.NewOrm().LoadRelated(bbs, "Images")

	var review []*models.BbsReview
	o.QueryTable(new(models.BbsReview)).RelatedSel().Filter("Status", 1).Filter("bbs_id", bbs_id).All(&review)
	bbs.BbsReview = review
	if ctl.IsLogin {
		bbs.IsLike = o.QueryTable(new(models.BbsLike)).Filter("bbs_id", bbs.Id).Filter("customer_id", ctl.CustomerId).Exist()
	} else {
		bbs.IsLike = false
	}

	if bbs.TopicId != 0 {
		var topic models.Topic
		err = o.QueryTable(new(models.Topic)).Filter("id", bbs.TopicId).Filter("status", 1).One(&topic)
		bbs.Topic = &topic
	}
	if bbs.ArticleId != 0 {
		var article models.Article
		err = o.QueryTable(new(models.Article)).Filter("id", bbs.ArticleId).Filter("status", 1).One(&article)
		if article.Cover == "" {
			img := utils.FindImg(article.Html)
			if len(img) > 0 {
				article.Cover = img[0]
			} else {
				article.Cover = "/static/images/bitimage.png"
			}

		}
		bbs.Article = &article
	}
	if len(bbs.Images) > 0 {
		for _, vv := range bbs.Images {
			ext := path.Ext(vv.Url)
			if ext == ".mp4" || ext == ".avi" || ext == ".mov" {
				vv.IsVideo = true
			}
		}
	}

	// var review = make(map[int][]utils.ReviewTree)

	ctl.Data["Review"] = map[int][]utils.ReviewTree{
		bbs_id: utils.ReviewTreeR(bbs.BbsReview, 0, 0, &models.BbsReview{}),
	}

	ctl.Data["Data"] = []*models.Bbs{
		bbs,
	}

	// 主题
	var topic []*admin.Topic
	o.QueryTable(new(admin.Topic)).Filter("status", 1).OrderBy("-join").Limit(10).All(&topic)
	ctl.Data["Topic"] = topic

	// 活跃用户
	var hotUser []*models.Customer
	o.QueryTable(new(models.Customer)).Filter("Email__isnull", false).OrderBy("-integral").Limit(10).Offset(0).All(&hotUser)
	for _, v := range hotUser {
		if ctl.IsLogin {
			v.IsFans = models.IsFans(ctl.CustomerId, v.Id)
		} else {
			v.IsFans = false
		}
	}

	ctl.Data["hotUser"] = hotUser
	if ctl.IsLogin {
		// 动态数量
		bbsCount, _ := o.QueryTable(new(models.Bbs)).Filter("status", 1).Filter("customer_id", ctl.CustomerId).Count()
		ctl.Data["BbsCount"] = bbsCount
		// 访问量
		ctl.Data["Visit"] = models.GetVisitByCustomerId(ctl.CustomerId)
	}
	// ctl.Data["Data"] = &bbs
	// ctl.Data["Paginator"] = utils.GenPaginator(page, limit, count)
	// ctl.Data["json"] = response
	f := func(source string, l int) string {
		var r = []rune(source)
		length := len(r)
		if l >= length {
			return source
		}

		return string(r[0:(l - 4)])
	}
	ctl.Data["SeoTitle"] = f(bbs.Content, 30)
	ctl.Data["SeoDescription"] = f(bbs.Content, 150)
	ctl.Data["IsTopic"] = false
	// 评论框是否显示
	ctl.Data["ReviewShow"] = true
	ctl.TplName = "home/" + ctl.Template + "/bbs_detail.html"
}

// http://v3.wufazhuce.com:8000/api/channel/one/0/0

func (c *BbsController) SaveDaily() {

	response := make(map[string]interface{})
	response["msg"] = "自动签到功能暂时关闭"
	response["code"] = 500
	response["err"] = "自动签到功能暂时关闭！"
	c.Data["json"] = response
	c.ServeJSON()

	c.StopRun()
	if !c.IsLogin {
		response["msg"] = "请先登录！"
		response["code"] = 500
		response["err"] = "请先登录！"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	content := c.GetString("content")

	content = strings.Trim(content, " ")
	dd, _ := utils.GetDailyData()
	if content == "" {
		content = "#每日签到# " + dd.Data.DailyContent[0].Forward + "——" + dd.Data.DailyContent[0].WordsInfo
	}

	//c.CustomerId = 653
	//c.Username = "我是甜美的西红柿"
	var topic string
	topic, content = utils.GetTopic(content)
	if topic != "#每日签到#" {
		response["msg"] = "非法操作"
		response["code"] = 500
		response["err"] = "非法操作"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	o := orm.NewOrm()
	err := o.Begin()

	bbs := models.Bbs{
		Customer: &models.Customer{Id: c.CustomerId, Username: c.Username},
		Content:  template.HTMLEscapeString(content),
		Created:  time.Now(),
		Updated:  time.Now(),
		Status:   1,
	}

	var topic_id int64

	if len(topic) != 0 {
		tp := models.Topic{Content: topic}

		err := o.Read(&tp, "Content")

		if err == orm.ErrNoRows {
			topic_id, err = o.Insert(&models.Topic{
				Content: topic,
				Created: time.Now(),
				Join:    1,
				Status:  1,
			})
			if err != nil {
				o.Rollback()
				response["msg"] = err.Error()
				response["code"] = 500
				response["err"] = err.Error()
				c.Data["json"] = response
				c.ServeJSON()
				c.StopRun()
			}
		} else if err == orm.ErrMissPK {
			o.Rollback()
			response["msg"] = err.Error()
			response["code"] = 500
			response["err"] = err.Error()
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		} else {
			tp.Join += 1
			if _, err := o.Update(&tp); err != nil {
				o.Rollback()
				response["msg"] = err.Error()
				response["code"] = 500
				response["err"] = err.Error()
				c.Data["json"] = response
				c.ServeJSON()
				c.StopRun()
			}
			topic_id = int64(tp.Id)
		}

		bbs.TopicId = int(topic_id)
	}

	var id int64
	id, err = o.Insert(&bbs)
	if err != nil {
		response["msg"] = "新增失败！"
		response["code"] = 500
		response["err"] = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	images := []models.BbsImages{}
	//img := c.GetStrings("bbsimg[]")
	img := []string{
		dd.Data.DailyContent[0].ImgUrl,
	}
	if len(img) > 0 {
		for _, v := range img {
			images = append(images, models.BbsImages{
				Bbs: &models.Bbs{Id: int(id)},
				Url: v,
			})
		}
		_, err := o.InsertMulti(len(images), images)
		if err != nil {
			o.Rollback()
			response["msg"] = "发布失败！"
			response["code"] = 500
			response["err"] = err.Error()
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		}

	}

	o.Commit()
	admin.UpdateIntegral(c.CustomerId, 5)
	response["msg"] = "发布成功！"
	response["code"] = 200
	response["id"] = id
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}
func (c *BbsController) Save() {

	response := make(map[string]interface{})
	content := c.GetString("content")
	//fmt.Printf("%v\n",content)
	img, content := utils.Find9Img(content)
	// fmt.Printf("%v,%v\n",img,ct)
	topic := c.GetString("topic")
	topic_content := c.GetString("topic_content")
	//fmt.Printf("%v\n",content)
	content = strings.Trim(content, " ")
	content = strings.Replace(content, "\n", "<br>", -1)
	if topic == "" {
		topic = topic_content
	}
	if content == "" && len(img) == 0 {
		response["msg"] = "内容不能为空！"
		response["code"] = 500
		response["err"] = "内容不能为空"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}
	if content == "" {
		content = "分享了图片"
	}
	if !c.IsLogin {
		response["msg"] = "请先登录！"
		response["code"] = 500
		response["err"] = "请先登录！"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}
	//c.CustomerId = 653
	//c.Username = "我是甜美的西红柿"

	cid, _ := c.GetInt64("cid", 0)

	o := orm.NewOrm()
	err := o.Begin()
	bbs := models.Bbs{
		Customer:   &models.Customer{Id: c.CustomerId, Username: c.Username},
		Content:    content,
		Created:    time.Now(),
		Updated:    time.Now(),
		Status:     1,
		CategoryId: int(cid),
	}

	var topic_id int64

	if len(topic) != 0 {
		topic = utils.GetTopicContent(topic)

		tp := models.Topic{Content: utils.GetTopicContent(topic)}

		err := o.Read(&tp, "Content")

		if err == orm.ErrNoRows {
			if len(img) > 0 {
				topic_id, err = o.Insert(&models.Topic{
					Content: topic,
					Created: time.Now(),
					Join:    1,
					Status:  1,
					Image:   img[0],
				})
			} else {
				topic_id, err = o.Insert(&models.Topic{
					Content: topic,
					Created: time.Now(),
					Join:    1,
					Status:  1,
				})
			}
			if err != nil {
				o.Rollback()
				response["msg"] = err.Error()
				response["code"] = 500
				response["err"] = err.Error()
				c.Data["json"] = response
				c.ServeJSON()
				c.StopRun()
			}
		} else if err == orm.ErrMissPK {
			o.Rollback()
			response["msg"] = err.Error()
			response["code"] = 500
			response["err"] = err.Error()
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		} else {
			tp.Join += 1
			if tp.Image == "" && len(img) > 0 {
				tp.Image = img[0]
			}
			if _, err := o.Update(&tp); err != nil {
				o.Rollback()
				response["msg"] = err.Error()
				response["code"] = 500
				response["err"] = err.Error()
				c.Data["json"] = response
				c.ServeJSON()
				c.StopRun()
			}
			topic_id = int64(tp.Id)
		}

		bbs.TopicId = int(topic_id)
	}

	var id int64
	id, err = o.Insert(&bbs)
	if err != nil {
		response["msg"] = "新增失败！"
		response["code"] = 500
		response["err"] = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	images := []models.BbsImages{}

	if len(img) > 0 {
		for _, v := range img {
			images = append(images, models.BbsImages{
				Bbs: &models.Bbs{Id: int(id)},
				Url: v,
			})
		}
		_, err := o.InsertMulti(len(images), images)
		if err != nil {
			o.Rollback()
			response["msg"] = "发布失败！"
			response["code"] = 500
			response["err"] = err.Error()
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		}

	}

	o.Commit()
	admin.UpdateIntegral(c.CustomerId, 5)

	// 通知关注用户
	focus := models.GetFocus(c.CustomerId)
	for _, v := range focus {
		utils.SendNotic(utils.NoticeMessage{
			Title:     "发了帖子",
			SendId:    c.CustomerId,
			ReceiveId: v.Id,
			Username:  c.Username,
			UserUrl:   fmt.Sprintf(`/user/%d`, c.CustomerId),
			Content:   bbs.Content,
			Replay:    "",
			Url:       fmt.Sprintf(`/bbs/detail/%d.html`, int(id)),
			Date:      time.Now(),
			Cover:     c.Cover,
		})
	}
	response["msg"] = "发布成功！"
	response["code"] = 200
	response["id"] = id
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *BbsController) SaveReview() {

	response := make(map[string]interface{})
	content := c.GetString("content")
	bbs_id, _ := c.GetInt("bbs_id", 0)
	reply_id, _ := c.GetInt("reply_id", 0)
	if content == "" || bbs_id == 0 {
		response["msg"] = "内容不能为空！"
		response["code"] = 500
		response["err"] = "内容不能为空"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}
	if !c.IsLogin {
		response["msg"] = "请先登录！"
		response["code"] = 500
		response["err"] = "请先登录！"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	BbsReview := models.BbsReview{
		Customer: &models.Customer{Id: c.CustomerId},
		Content:  template.HTMLEscapeString(content),
		Created:  time.Now(),
		Status:   1,
		Bbs:      &models.Bbs{Id: bbs_id},
		ReplyId:  reply_id,
	}
	o := orm.NewOrm()
	err := o.Begin()
	var id int64
	id, err = o.Insert(&BbsReview)
	if err != nil {
		response["msg"] = "新增失败！"
		response["code"] = 500
		response["err"] = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	bbs := admin.Bbs{Id: bbs_id}
	if o.Read(&bbs) == nil {
		bbs.Review += 1
		bbs.Updated = time.Now()
		if _, err := o.Update(&bbs); err != nil {
			o.Rollback()
			response["msg"] = "修改失败！"
			response["code"] = 500
			response["err"] = err.Error()
		}
	}

	o.Commit()
	admin.UpdateIntegral(c.CustomerId, 2)
	var receiveId int
	if reply_id == 0 {
		receiveId = bbs.Customer.Id
		utils.SendNotic(utils.NoticeMessage{
			Title:     "评论了你的帖子",
			SendId:    c.CustomerId,
			ReceiveId: receiveId,
			Username:  c.Username,
			UserUrl:   fmt.Sprintf(`/user/%d`, c.CustomerId),
			Content:   bbs.Content,
			Replay:    content,
			Url:       fmt.Sprintf(`/bbs/detail/%d.html`, bbs.Id),
			Date:      time.Now(),
			Cover:     c.Cover,
		})
	} else {
		brmodel, err := admin.GetBbsReviewById(reply_id)
		if err == nil {
			receiveId = brmodel.Customer.Id
		}
		utils.SendNotic(utils.NoticeMessage{
			Title:     "评论了你的评论",
			SendId:    c.CustomerId,
			ReceiveId: receiveId,
			Username:  c.Username,
			UserUrl:   fmt.Sprintf(`/user/%d`, c.CustomerId),
			Content:   brmodel.Content,
			Replay:    content,
			Url:       fmt.Sprintf(`/bbs/detail/%d.html`, bbs.Id),
			Date:      time.Now(),
			Cover:     c.Cover,
		})
	}

	type EmailData struct {
		Review admin.BbsReview
		Bbs    admin.Bbs
		User   string
	}

	// 发送邮件
	customerEmail, err := admin.GetCustomerById(int(receiveId))
	if customerEmail != nil && customerEmail.Email != "" && err != nil {
		flag := utils.SendEmail(utils.Email{
			From: "1920853199@qq.com",
			To:   customerEmail.Email,
			Header: map[string]string{
				"Subject": c.Username + "评论了您的帖子",
			},
			Template: "views/home/nihongdengxia/bbs-comment-email.html",
			Data: EmailData{
				Review: BbsReview,
				Bbs:    bbs,
				User:   c.Username,
			},
		})
		if flag == nil {
			fmt.Println("评论邮件发送成功！")
		} else {
			fmt.Println("评论邮件发送失败！")
		}
	}

	response["msg"] = "评论成功！"
	response["code"] = 200
	response["id"] = id
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()

}

func (c *BbsController) Like() {

	bbs_id, _ := c.GetInt("bbs_id", 0)
	response := make(map[string]interface{})
	if bbs_id == 0 {
		response["msg"] = "内容不能为空！"
		response["code"] = 500
		response["err"] = "内容不能为空"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	if !c.IsLogin {
		response["msg"] = "请先登录！"
		response["code"] = 500
		response["err"] = "请先登录！"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}
	o := orm.NewOrm()
	var isAdd = 1
	flag := o.QueryTable(new(models.BbsLike)).Filter("bbs_id", bbs_id).Filter("customer_id", c.CustomerId).Exist()
	if flag {
		isAdd = 2
	}

	err := o.Begin()
	if isAdd == 2 {
		_, err := o.QueryTable(new(models.BbsLike)).Filter("bbs_id", bbs_id).Filter("customer_id", c.CustomerId).Delete()
		if err != nil {
			o.Rollback()
			response["msg"] = "操作失败！"
			response["code"] = 500
			response["err"] = "操作失败！"
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		}
	} else {
		bbslike := models.BbsLike{
			BbsId:      bbs_id,
			CustomerId: c.CustomerId,
		}
		_, err = o.Insert(&bbslike)
		if err != nil {
			o.Rollback()
			response["msg"] = "点赞失败！"
			response["code"] = 500
			response["err"] = err.Error()
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		}
	}

	bbs := admin.Bbs{Id: bbs_id}

	if o.Read(&bbs) == nil {
		if isAdd == 1 {
			bbs.Like += 1
		} else {
			bbs.Like -= 1
		}
		if _, err := o.Update(&bbs); err != nil {
			o.Rollback()
			response["msg"] = "点赞失败！"
			response["code"] = 500
			response["err"] = err.Error()
		} else {
			o.Commit()
			admin.UpdateIntegral(c.CustomerId, 2)
			if isAdd == 1 {
				utils.SendNotic(utils.NoticeMessage{
					Title:     "赞了你的帖子",
					SendId:    c.CustomerId,
					ReceiveId: bbs.Customer.Id,
					Username:  c.Username,
					UserUrl:   fmt.Sprintf(`/user/%d`, c.CustomerId),
					Content:   bbs.Content,
					Replay:    "",
					Url:       fmt.Sprintf(`/bbs/detail/%d.html`, bbs.Id),
					Date:      time.Now(),
					Cover:     c.Cover,
				})
			}
			response["msg"] = "点赞成功"
			response["code"] = 200
			response["review"] = bbs.Like
			response["type"] = isAdd
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		}
	}
}

func (ctl *BbsController) Topic() {
	// if !ctl.IsLogin {
	// 	ctl.Redirect("/login.html", 302)
	// }

	limit := int64(50)
	page, _ := ctl.GetInt64("page", 1) // 页数
	offset := (page - 1) * limit       // 偏移量
	//topic_id, _ := ctl.GetInt64("topic")
	id := ctl.Ctx.Input.Param(":id")
	topic_id, _ := strconv.Atoi(id)

	o := orm.NewOrm()

	// topic_id = 1
	if topic_id < 1 {
		ctl.Abort("404")
	}
	var topicContent models.Topic
	err := o.QueryTable(new(models.Topic)).Filter("id", topic_id).Filter("status", 1).One(&topicContent)
	if err != nil {
		ctl.Abort("404")
	}
	ctl.Data["TopicContent"] = &topicContent

	var bbs []*models.Bbs
	qs := o.QueryTable(new(models.Bbs))
	qs = qs.Filter("status", 1)
	qs = qs.Filter("Customer__Username__isnull", false)
	qs = qs.Filter("topic_id", topic_id)

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

		if !ctl.IsLogin {
			v.IsLike = false
		} else {
			v.IsLike = o.QueryTable(new(models.BbsLike)).Filter("bbs_id", v.Id).Filter("customer_id", ctl.CustomerId).Exist()
		}

	}
	var review = make(map[int][]utils.ReviewTree)
	for _, v := range bbs {
		review[v.Id] = utils.ReviewTreeR(v.BbsReview, 0, 0, &models.BbsReview{})
	}

	ctl.Data["Review"] = &review
	ctl.Data["Data"] = &bbs
	ctl.Data["Paginator"] = utils.GenPaginator(page, limit, count)
	ctl.Data["ReviewShow"] = false
	// 主题
	var topic []*admin.Topic
	o.QueryTable(new(admin.Topic)).Filter("status", 1).OrderBy("-join").Limit(10).All(&topic)
	ctl.Data["Topic"] = topic

	// 活跃用户
	var hotUser []*models.Customer
	o.QueryTable(new(models.Customer)).Filter("Email__isnull", false).OrderBy("-integral").Limit(5).Offset(0).All(&hotUser)
	for _, v := range hotUser {
		if !ctl.IsLogin {
			v.IsFans = false
		} else {
			v.IsFans = models.IsFans(ctl.CustomerId, v.Id)
		}
	}

	ctl.Data["hotUser"] = hotUser
	// ctl.Data["json"] = bbs
	// ctl.ServeJSON()
	// ctl.StopRun()
	ctl.Data["IsTopic"] = true
	ctl.TplName = "home/" + ctl.Template + "/topic.html"
}

func (ctl *BbsController) SaveView() {

	view := ctl.GetString(":splat")
	if view == "program" || view == "bbs" {
		ctl.Ctx.SetCookie("bbsview", view)

	}
	ctl.Redirect("/bbs.html", 302)
}
