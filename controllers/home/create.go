package home

import (
	"fmt"
	"go-bbs/models"
	"go-bbs/models/admin"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type CustomerArticleController struct {
	BaseController
}

func (c *CustomerArticleController) Create() {
	if !c.IsLogin {
		c.Redirect("/login.html"+"?referrer="+c.Ctx.Input.URL(), 302)
		//c.Redirect("/login.html", 302)
	}
	set, _ := models.GetArticleSetCustomerId(c.CustomerId)

	c.Data["Set"] = set
	c.TplName = "home/" + c.Template + "/articles_create.html"
}

func (c *CustomerArticleController) Edit() {

	response := make(map[string]interface{})

	if c.IsLogin == false {
		response["msg"] = "请先登录！"
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	id, _ := c.GetInt("id")
	title := c.GetString("title")
	tag := c.GetString("tag")
	cate, _ := c.GetInt("cate", 0)
	remark := c.GetString("remark")
	desc := c.GetString("desc_content")
	html := c.GetString("desc_html")

	o := orm.NewOrm()

	article := admin.Article{Id: id}
	if o.Read(&article) == nil {
		if article.Customer.Id != c.CustomerId {
			response["msg"] = "非法操作！"
			response["code"] = 500
			response["err"] = "非法操作！"
		}
		article.Title = title
		article.Tag = tag
		article.Desc = desc
		article.Html = html
		article.Remark = remark

		valid := validation.Validation{}
		valid.Required(article.Title, "Title")
		valid.Required(article.Tag, "Tag")
		valid.Required(article.Desc, "Desc")
		valid.Required(article.Html, "Html")
		//valid.Required(article.Remark, "Remark")

		if valid.HasErrors() {
			// 如果有错误信息，证明验证没通过
			// 打印错误信息
			for _, err := range valid.Errors {
				//log.Println(err.Key, err.Message)
				response["msg"] = "修改失败！"
				response["code"] = 500
				response["err"] = err.Key + " " + err.Message
				c.Data["json"] = response
				c.ServeJSON()
				c.StopRun()
			}
		}

		article.Category = &admin.Category{cate, "", 0, 0, 0, ""}

		if _, err := o.Update(&article); err == nil {
			response["msg"] = "修改成功！"
			response["code"] = 200
			response["id"] = id
		} else {
			response["msg"] = "修改失败！"
			response["code"] = 500
			response["err"] = err.Error()
		}
	} else {
		response["msg"] = "修改失败！"
		response["code"] = 500
		response["err"] = "ID 不能为空！"
	}
	SaveSet(c, article)
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *CustomerArticleController) Put() {
	if !c.IsLogin {
		c.Redirect("/login.html"+"?referrer="+c.Ctx.Input.URL(), 302)
		//c.Redirect("/login.html", 302)
	}

	id, err := c.GetInt("id", 0)

	if id == 0 {
		c.Abort("404")
	}

	set, _ := models.GetArticleSetCustomerId(c.CustomerId)

	c.Data["Set"] = set

	// 基础数据
	o := orm.NewOrm()
	var articles admin.Article
	qs := o.QueryTable(new(admin.Article))
	err = qs.Filter("id", id).Filter("customer_id", c.CustomerId).One(&articles)
	if err != nil {
		c.Abort("404")
	}
	c.Data["Data"] = &articles

	// fmt.Println(tempArticles)
	setId, err := models.GetArticleSetByArticleId(articles.Id, c.CustomerId)
	fmt.Println(setId)
	if err == nil {
		c.Data["SetId"] = setId.Id
	} else {
		c.Data["SetId"] = 0
	}
	c.TplName = "home/" + c.Template + "/articles_edit.html"
}

func (c *CustomerArticleController) Save() {
	response := make(map[string]interface{})

	if c.IsLogin == false {
		response["msg"] = "请先登录！"
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	title := c.GetString("title")
	tag := c.GetString("tag")
	cate, _ := c.GetInt("cate", 0)
	remark := c.GetString("remark")
	desc := c.GetString("desc_content")
	html := c.GetString("desc_html")
	// url := c.GetString("url")
	// cover := c.GetString("cover")

	o := orm.NewOrm()
	article := admin.Article{
		Title:    title,
		Tag:      tag,
		Desc:     desc,
		Html:     html,
		Remark:   remark,
		Status:   1,
		Customer: &admin.Customer{Id: c.CustomerId},
		Category: &admin.Category{cate, "", 0, 0, 0, ""},
	}

	valid := validation.Validation{}
	valid.Required(article.Title, "Title")
	valid.Required(article.Category, "Category")
	valid.Required(article.Html, "Html")
	//valid.Required(article.Tag, "Tag")
	valid.Required(article.Desc, "Desc")
	//valid.Required(article.Remark, "Remark")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			//log.Println(err.Key, err.Message)
			response["msg"] = "新增失败！"
			response["code"] = 500
			response["err"] = err.Key + " " + err.Message
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		}
	}

	if id, err := o.Insert(&article); err == nil {
		response["msg"] = "新增成功！"
		response["code"] = 200
		response["id"] = id
	} else {
		response["msg"] = "新增失败！"
		response["code"] = 500
		response["err"] = err.Error()
	}

	SaveSet(c, article)
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()

}

func SaveSet(c *CustomerArticleController, a admin.Article) {
	set := c.GetString("set")
	set_content := c.GetString("set_content")
	o := orm.NewOrm()
	if set == "" {
		set = set_content
	}
	if set == "" {
		return
	}

	var l models.ArticleSet
	qs := o.QueryTable(new(models.ArticleSet))
	qs = qs.Filter("customer_id", c.CustomerId)
	qs = qs.Filter("status", 1)
	qs = qs.Filter("article_id", a.Id)
	//qs = qs.Filter("title", a.Id)
	// v = &ArticleSet{ArticleId: ArticleId, Status: 1, CustomerId: CustomerId}
	err := qs.One(&l)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(l)
	}

	if err == orm.ErrNoRows {
		models.AddArticleSet(&models.ArticleSet{
			ArticleId:  a.Id,
			CustomerId: c.CustomerId,
			Created:    time.Now(),
			Status:     1,
			Title:      set,
		})

	} else if err == orm.ErrMissPK {

	} else {
		models.UpdateArticleSetById(&models.ArticleSet{
			Id:         l.Id,
			ArticleId:  a.Id,
			CustomerId: c.CustomerId,
			Title:      set,
			Created:    l.Created,
			Status:     l.Status,
		})
	}
}
