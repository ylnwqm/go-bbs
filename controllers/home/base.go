package home

import (
	"fmt"
	"go-bbs/models/admin"
	"go-bbs/utils"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type BaseController struct {
	beego.Controller
	Template   string
	IsLogin    bool
	CustomerId int
	Username   string
	Cover      string
}

func (c *BaseController) Layout() {

	o := orm.NewOrm()

	category := new(admin.Category)

	var categorys []*admin.Category
	qs := o.QueryTable(category)
	qs = qs.Filter("status", 1)

	qs.OrderBy("-sort").All(&categorys)

	tree := utils.CategoryTree(categorys, 0, 0)
	// c.Data["json"] = tree
	// c.ServeJSON()
	// c.StopRun()

	c.Data["Category"] = tree

	// 月份排序
	articleTime := new(admin.Article)
	var articlesTime []*admin.Article
	nqs := o.QueryTable(articleTime)
	nqs = nqs.Filter("status", 1)
	nqs.OrderBy("-Created").RelatedSel().All(&articlesTime, "Created")
	count, _ := nqs.Count()
	var datetime = make(map[string]int64)
	var dateTimeKey []string
	for _, v := range articlesTime {
		//str = append(str ,v.Created.Format("2006-01"))
		//c.Ctx.WriteString(v.Created.Format("2006-01"))
		k := v.Created.Format("2006-01")
		if datetime[k] == 0 {
			dateTimeKey = append(dateTimeKey, k)
		}
		datetime[k] = datetime[k] + 1
	}
	c.Data["DateTime"] = datetime
	c.Data["DateTimeKey"] = dateTimeKey
	c.Data["DateCount"] = count

	// 阅读排序
	articleReadSort := new(admin.Article)
	var articlesReadSort []*admin.Article
	nqrs := o.QueryTable(articleReadSort)
	nqrs = nqrs.Filter("status", 1)
	nqrs = nqrs.OrderBy("-Pv")
	nqrs.Limit(10).All(&articlesReadSort, "Id", "Title", "Pv")
	c.Data["ArticlesReadSort"] = articlesReadSort

	// 推荐
	articleRecommend := new(admin.Article)
	var articlesRecommend []*admin.Article
	nqr := o.QueryTable(articleRecommend)
	nqr = nqr.Filter("status", 1)
	nqr = nqr.Filter("recommend", 1)
	nqr = nqr.OrderBy("-Id")
	nqr.Limit(5).All(&articlesRecommend, "Id", "Title", "Pv", "Cover")
	c.Data["ArticlesRecommend"] = articlesRecommend

	// 最新排序
	var articlesNewSort []*admin.Article
	nqns := o.QueryTable(new(admin.Article))
	nqns = nqns.Filter("status", 1)
	nqns = nqns.OrderBy("-Id")
	nqns.Limit(15).All(&articlesNewSort, "Id", "Title", "Pv", "Created")
	c.Data["ArticlesNewSort"] = articlesNewSort

	// 最新评论
	review := new(admin.Review)
	var reviewData []*admin.Review
	nqrw := o.QueryTable(review)
	nqrw = nqrw.Filter("status", 1)
	nqrw = nqrw.OrderBy("-Id")
	nqrw.Limit(5).All(&reviewData, "Review", "ArticleId")
	reviewCount, _ := nqrw.Count()
	c.Data["ReviewCount"] = reviewCount
	c.Data["Review"] = reviewData

	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10000
	var offset int64

	sortby = append(sortby, "sort")
	order = append(order, "asc")
	link, _ := admin.GetAllLink(query, fields, sortby, order, offset, limit)
	c.Data["Link"] = link
}

func (c *BaseController) Menu() {

	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	sortby = append(sortby, "sort")
	order = append(order, "asc")

	menu, _ := admin.GetAllMenu(query, fields, sortby, order, offset, limit)
	data := utils.MenuData(menu, 0, 0)
	/*c.Data["json"] = data
	c.ServeJSON()
	c.StopRun()*/
	c.Data["Menu"] = data

}

var mobileRe, _ = regexp.Compile("(?i:Mobile|iPod|iPhone|Android|Opera Mini|BlackBerry|webOS|UCWEB|Blazer|PSP)")

func (c *BaseController) Prepare() {

	// fmt.Printf("前一个页面：%s\n", c.Ctx.Request.RequestURI)
	// 登录信息
	customer := c.GetSession("Customer")
	//fmt.Println("登录信息：", customer)
	if customer != nil {
		//ctl.User = customer.(admin.Customer)
		c.Data["Customer"] = customer.(admin.Customer)
		//c.Data["CustomerId"] = customer.(admin.Customer).Id
		c.Data["IsLogin"] = true
		c.IsLogin = true
		c.CustomerId = customer.(admin.Customer).Id
		c.Username = customer.(admin.Customer).Username
		c.Cover = customer.(admin.Customer).Image
		// notice
		noticeCount, _ := admin.GetNoticeCount(c.CustomerId)
		c.Data["NoticeCount"] = noticeCount

	} else {
		c.Data["IsLogin"] = false
		c.IsLogin = false
	}

	c.Data["bgClass"] = "bgColor"
	c.Data["T"] = time.Now()
	c.Data["WT"] = utils.WeekDayMap[time.Now().Weekday().String()]

	// 设置信息
	o := orm.NewOrm()
	var setting []*admin.Setting
	o.QueryTable(new(admin.Setting)).All(&setting)

	for _, v := range setting {
		c.Data[v.Name] = v.Value
		// if v.Name == "template" {
		// 	 c.Template = v.Value
		// }
	}

	ua := mobileRe.FindString(c.Ctx.Input.UserAgent())
	if ua == "" {
		c.Template = c.Data["template"].(string)
	} else {
		c.Template = c.Data["mobile"].(string)
	}

	// fmt.Printf("Host：%s\n",c.Ctx.Input.IP())
	if c.Ctx.Input.Host() == "clblog.club" {
		c.Template = "toutiao"
	}

	if c.Ctx.Input.Host() == "gooooooooogle.cn" || c.Ctx.Input.Host() == "www.gooooooooogle.cn" {
		c.Template = "baidu"
	}

	// 站点统计信息
	// pv, _ := o.QueryTable(new(admin.Log)).Count()
	// uv, _ := o.QueryTable(new(admin.Log)).GroupBy("ip").Count()

	c.Data["PV"] = "暂时关闭PV统计展示"
	c.Data["UV"] = "暂时关闭UV统计展示"
	// todayUv, _ := o.QueryTable(new(admin.Log)).Filter("create__gte", time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())).GroupBy("ip").Count()
	c.Data["TUV"] = "关闭显示"
	// Ad
	var fields []string
	var sortby = []string{
		"id",
	}
	var order = []string{
		"desc",
	}
	var query = make(map[string]string)
	// Ad
	l, err := admin.GetAllAd(query, fields, sortby, order, 0, 0)
	if err != nil {
		c.Abort(err.Error())
	}
	type Data struct {
		Gid   string
		Group string
		Ad    []*admin.Ad
	}

	var data []*Data

	for _, value := range l {
		ad := value.(admin.Ad)
		var flag bool = false
		for _, v := range data {
			if v.Gid == ad.Gid {
				v.Ad = append(v.Ad, &ad)
				flag = true
				break
			}
		}

		if flag == true {
			continue
		}

		data = append(data, &Data{
			Gid:   ad.Gid,
			Group: ad.Group,
			Ad:    []*admin.Ad{&ad},
		})
	}

	for _, v := range data {
		c.Data[v.Gid] = v.Ad
	}

	// music
	// type m struct {
	// 	Name   string `json:"name"`
	// 	Url    string `json:"url"`
	// 	Cover  string `json:"cover"`
	// 	Author string `json:"author"`
	// }
	// var music []m
	// var lists []orm.ParamsList
	// num, err := o.Raw("SELECT * FROM music order by RAND() limit 10").ValuesList(&lists)
	// //num, err := o.Raw("SELECT * FROM music order by id asc limit 10").ValuesList(&lists)
	// if err == nil && num > 0 {
	// 	for _, v := range lists {
	// 		music = append(music, m{
	// 			//Id:     v[0].(int),
	// 			Name:   v[1].(string),
	// 			Url:    v[2].(string),
	// 			Cover:  v[3].(string),
	// 			Author: v[4].(string),
	// 			//SongId: v[5].(string),
	// 		})
	// 	}
	// }

	// if len(music) == 0 {
	// 	c.Data["json"] = "[]"
	// } else {
	// 	c.Data["json"] = music
	// }
	//mu, err := json.Marshal(music)
	//c.Data["Music"] = string(mu)
	//c.Log()

	c.Layout()
	c.Menu()
	c.Keywords()
	c.Log()
}

func (c *BaseController) Keywords() {

	o := orm.NewOrm()
	qs := o.QueryTable(new(admin.Article))

	var tag []*admin.Article

	qs = qs.Filter("status", 1)
	qs = qs.Filter("Customer__Username__isnull", false)
	qs = qs.Filter("Category__Name__isnull", false)
	qs.All(&tag, "tag")

	var tags []string
	for _, v := range tag {
		tags = append(tags, strings.Split(strings.Replace(v.Tag, `，`, `,`, -1), `,`)...)
	}

	var tagsMap = make(map[string]int)

	for _, v := range tags {
		tagsMap[v] += 1
	}

	for k, _ := range tagsMap {
		tagsMap[k] = tagsMap[k]/5 + 15
	}

	tagsMap = map[string]int{"资源下载": 0, "社区": 0, "霓虹灯下": 0, "读书": 0, "书评": 0, "新鲜事": 0, "摄影": 0, "有趣": 0, "好玩": 0, "文字": 0, "休闲": 0, "故事": 0, "技术编程": 0, "代码": 0, "有个性的网站": 0, "热点动态": 0}
	c.Data["Tag"] = tagsMap
	c.Data["KTag"] = tagsMap

}

func (c *BaseController) Log() {

	ua := mobileRe.FindString(c.Ctx.Input.UserAgent())

	var uaStr string
	if ua == "" {
		uaStr = "PC端"
	} else {
		uaStr = "移动端"
	}

	path := c.Ctx.Input.URL()
	if c.Ctx.Input.URL() == "/article/like" {
		id, _ := c.GetInt("id")
		path = c.Ctx.Input.URL() + strconv.Itoa(id)
	}
	o := orm.NewOrm()
	var log = admin.Log{
		Ip:        c.Ctx.Input.IP(),
		City:      utils.IP(c.Ctx.Input.IP()),
		UserAgent: uaStr,
		Page:      path,
		Referer:   c.Ctx.Input.Referer(),
	}
	_, err := o.Insert(&log)
	if err != nil {
		fmt.Println(err.Error())
	}
}
