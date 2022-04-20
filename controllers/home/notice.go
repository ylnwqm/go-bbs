package home

import (
	"fmt"
	"go-bbs/models/admin"
	"html/template"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type NoticeController struct {
	BaseController
}

func (c *NoticeController) Notice() {
	if !c.IsLogin {
		c.Redirect("/login.html"+"?referrer="+c.Ctx.Input.URL(), 302)
	}
	var nav []map[string]string

	nav = append(nav, map[string]string{
		"Title": "首页",
		"Url":   "/",
	})
	nav = append(nav, map[string]string{
		"Title": "我的消息",
		"Url":   "/notice.html",
	})
	c.Data["Nav"] = nav
	c.Data["Navlen"] = len(nav) - 1

	n, err := admin.GetNiticeByReceiveId(c.CustomerId, 30, 0)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	c.Data["Notice"] = n
	c.TplName = "home/" + c.Template + "/notice.html"

}
func (c *NoticeController) Send() {
	if !c.IsLogin {
		c.Redirect("/login.html"+"?referrer"+c.Ctx.Input.URL(), 302)
	}

	uid, _ := c.GetInt("uid")
	msg := c.GetString("msg")
	o := orm.NewOrm()
	notice := admin.Notice{
		SendId:    c.CustomerId,
		ReceiveId: uid,
		Type:      admin.MsgMessage,
		Title:     admin.MsgTitle[admin.MsgMessage],
		Content:   template.HTMLEscapeString(msg),
		Status:    1,
	}

	response := make(map[string]interface{})

	valid := validation.Validation{}

	valid.Required(notice.ReceiveId, "ReceiveId")
	valid.Required(notice.Content, "Content")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			//log.Println(err.Key, err.Message)
			response["msg"] = "发送失败"
			response["code"] = 500
			response["err"] = err.Key + " " + err.Message
			c.Data["json"] = response
			c.ServeJSON()
			c.StopRun()
		}
	}

	if id, err := o.Insert(&notice); err == nil {
		response["msg"] = "发送成功"
		response["code"] = 200
		response["id"] = id
	} else {
		response["msg"] = "发送失败"
		response["code"] = 500
		response["err"] = err.Error()
	}

	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}
