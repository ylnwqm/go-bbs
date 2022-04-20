package home

import (
	"go-bbs/models/admin"
	models "go-bbs/models/admin"
	"go-bbs/utils"
	histoday "go-bbs/utils/hotnews"
	"strconv"
	"unsafe"

	"github.com/astaxie/beego/orm"

	// "github.com/astaxie/beego/cache"
	"fmt"
	"strings"
	"time"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	page := c.GetString("page")
	// 推荐
	o := orm.NewOrm()
	var list []*admin.Article
	o.QueryTable(new(admin.Article)).Filter("status", 1).Filter("recommend", 1).Filter("Customer__Username__isnull", false).Filter("Category__Name__isnull", false).OrderBy("-id").RelatedSel().All(&list, "id", "title")
	c.Data["Recommend"] = list

	c.Data["index"] = "首页"

	var topic []*admin.Topic
	o.QueryTable(new(admin.Topic)).Filter("status", 1).OrderBy("-join").Limit(5).All(&topic)
	c.Data["Topic"] = topic

	if c.Template == "app" || c.Template == "toutiao" {
		((*ArticleController)(unsafe.Pointer(c))).List()
	}

	// if c.Template == "baidu" {
	// 	((*HotnewsController)(unsafe.Pointer(c))).List()
	// }

	// 活跃用户
	var hotUser []*models.Customer
	o.QueryTable(new(models.Customer)).Filter("Email__isnull", false).OrderBy("-id").Limit(8).Offset(0).All(&hotUser)
	for k, v := range hotUser {
		hotUser[k].Name = utils.Substring(v.Username, 4)
	}
	c.Data["HotUser"] = hotUser

	var hotnews []*admin.Hotnews
	o.QueryTable(new(admin.Hotnews)).RelatedSel().OrderBy("-id").Limit(30).All(&hotnews)

	for _, v := range hotnews {
		//fmt.Println(v.Category)
		v.Rurl = v.Url
		if c.Template == "baidu" {
			v.Url = "/hotnews/detail/" + strconv.Itoa(v.Id) + ".html"
		} else {
			v.Url = "/dailyhot/detail/" + strconv.Itoa(v.Id) + ".html"
		}
	}

	c.Data["Hotnews"] = hotnews

	if c.Template == "baidu" {

		// 百度热榜
		var baiduhotnews []*admin.Hotnews
		o.QueryTable(new(admin.Hotnews)).Filter("category_id", 27).RelatedSel().OrderBy("-id").Limit(20).All(&baiduhotnews)

		for _, v := range baiduhotnews {
			//fmt.Println(v.Category)
			v.Rurl = v.Url
			v.Url = "/hotnews/detail/" + strconv.Itoa(v.Id) + ".html"
		}

		for i, j := 0, len(baiduhotnews)-1; i < j; i, j = i+1, j-1 {
			baiduhotnews[i], baiduhotnews[j] = baiduhotnews[j], baiduhotnews[i]
		}

		c.Data["Baiduhotnews"] = baiduhotnews

		// 微博热榜
		var weibohotnews []*admin.Hotnews
		o.QueryTable(new(admin.Hotnews)).Filter("category_id", 28).RelatedSel().OrderBy("-id").Limit(20).All(&weibohotnews)
		for _, v := range weibohotnews {
			//fmt.Println(v.Category)
			v.Rurl = v.Url
			v.Url = "/hotnews/detail/" + strconv.Itoa(v.Id) + ".html"
		}

		for i, j := 0, len(weibohotnews)-1; i < j; i, j = i+1, j-1 {
			weibohotnews[i], weibohotnews[j] = weibohotnews[j], weibohotnews[i]
		}
		c.Data["Weibohotnews"] = weibohotnews

		// 知乎热榜
		var zhihuhotnews []*admin.Hotnews
		o.QueryTable(new(admin.Hotnews)).Filter("category_id", 21).RelatedSel().OrderBy("-id").Limit(20).All(&zhihuhotnews)
		for _, v := range zhihuhotnews {
			//fmt.Println(v.Category)
			v.Rurl = v.Url
			v.Url = "/hotnews/detail/" + strconv.Itoa(v.Id) + ".html"
		}
		for i, j := 0, len(zhihuhotnews)-1; i < j; i, j = i+1, j-1 {
			zhihuhotnews[i], zhihuhotnews[j] = zhihuhotnews[j], zhihuhotnews[i]
		}
		c.Data["Zhihuhotnews"] = zhihuhotnews

		// b站热榜
		var bilibilihotnews []*admin.Hotnews
		o.QueryTable(new(admin.Hotnews)).Filter("category_id", 20).RelatedSel().OrderBy("-id").Limit(20).All(&bilibilihotnews)
		for _, v := range bilibilihotnews {
			//fmt.Println(v.Category)
			v.Rurl = v.Url
			v.Url = "/hotnews/detail/" + strconv.Itoa(v.Id) + ".html"
		}
		for i, j := 0, len(bilibilihotnews)-1; i < j; i, j = i+1, j-1 {
			bilibilihotnews[i], bilibilihotnews[j] = bilibilihotnews[j], bilibilihotnews[i]
		}
		c.Data["Bilibilihotnews"] = bilibilihotnews

		// bm, _ := cache.NewCache("file", `{"CachePath":"./cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"120"}`)

		// 历史上的今天
		h := histoday.Histoday{}
		Histoday := h.Get("https://hao.360.com/histoday/")
		firstHistoday := Histoday[0:6]
		moreHistoday := Histoday[6:]
		c.Data["Histoday"] = firstHistoday
		c.Data["MoreHistoday"] = moreHistoday
		// 球赛
		p := histoday.Hupu{}
		c.Data["Match"] = p.GetMatch("https://www.hupu.com/")
		//bm.Put("Match",c.Data["Match"],5*time.Minute)

		// match := bm.Get("Match")
		// fmt.Printf("%v\n",match)

		danmu := c.Ctx.GetCookie("SwitchDanmu")
		if danmu == "" {
			danmu = "on"
		}

		c.Data["SwitchDanmu"] = danmu
	}

	if c.Template == "nihongdengxia" {
		var bbs []*models.Bbs
		qs := o.QueryTable(new(models.Bbs))
		qs = qs.Filter("status", 1)

		qs = qs.Filter("Customer__Username__isnull", false)

		// 获取数据
		_, err := qs.OrderBy("-updated").RelatedSel().Limit(30).All(&bbs)
		if err != nil {
			panic(err)
		}

		for _, v := range bbs {

			if v.CategoryId != 0 {
				var category models.Category
				err = o.QueryTable(new(models.Category)).Filter("id", v.CategoryId).Filter("status", 1).One(&category)
				v.Category = &category
			}
			v.Content = strings.Replace(v.Content, "<br>", "", -1)
			v.Content = utils.BbsSubbbs(v.Content, 15)

		}

		c.Data["IndexBbs"] = bbs

		type Data struct {
			Image string
			Title string
			Url   string
			Date  time.Time
			Num   int
		}
		type Item struct {
			Name string
			Item []Data
		}
		var alldata []Item
		// 文章
		var articlesNewSort []*admin.Article
		nqns := o.QueryTable(new(admin.Article))
		nqns = nqns.Filter("status", 1)
		nqns = nqns.OrderBy("-Id")
		nqns.Limit(10).All(&articlesNewSort, "Id", "Title", "Pv", "Created", "Html")
		var arrArticlesData []Data
		for _, v := range articlesNewSort {
			var tempdata Data
			tempdata.Title = v.Title
			tempdata.Url = "/detail/" + fmt.Sprintf("%d", v.Id) + ".html"
			tempdata.Date = v.Created
			tempdata.Num = v.Pv
			if v.Cover == "" {
				img := utils.FindImg(v.Html)
				if len(img) > 0 {
					tempdata.Image = img[0]
				} else {
					tempdata.Image = "/static/images/bitimage.png"
				}
			} else {
				tempdata.Image = v.Cover
			}
			arrArticlesData = append(arrArticlesData, tempdata)

		}

		alldata = append(alldata, Item{
			Name: "文章",
			Item: arrArticlesData,
		})

		// 生活
		var arrBbsData12 []Data
		var bbs12 []*models.Bbs
		qs12 := o.QueryTable(new(models.Bbs))
		qs12 = qs12.Filter("status", 1)
		qs12 = qs12.Filter("category_id__in", 12)
		qs12 = qs12.Filter("Customer__Username__isnull", false)
		_, err = qs12.OrderBy("-id").RelatedSel().Limit(10).All(&bbs12)
		for _, v := range bbs12 {
			var tempdata Data
			tempdata.Title = v.Content
			tempdata.Url = "/bbs/detail/" + fmt.Sprintf("%d", v.Id) + ".html"
			tempdata.Date = v.Created
			tempdata.Num = v.Review
			orm.NewOrm().LoadRelated(v, "Images")
			if len(v.Images) > 0 {
				tempdata.Image = v.Images[0].Url
			} else {
				tempdata.Image = "/static/images/bitimage.png"
			}
			arrBbsData12 = append(arrBbsData12, tempdata)
		}

		alldata = append(alldata, Item{
			Name: "生活",
			Item: arrBbsData12,
		})

		// 美食
		var arrBbsData11 []Data
		var bbs11 []*models.Bbs
		qs = o.QueryTable(new(models.Bbs))
		qs = qs.Filter("status", 1)
		qs = qs.Filter("category_id__in", 11)
		qs = qs.Filter("Customer__Username__isnull", false)
		_, err = qs.OrderBy("-id").RelatedSel().Limit(10).All(&bbs11)
		for _, v := range bbs11 {
			var tempdata Data
			tempdata.Title = v.Content
			tempdata.Url = "/bbs/detail/" + fmt.Sprintf("%d", v.Id) + ".html"
			tempdata.Date = v.Created
			tempdata.Num = v.Review
			orm.NewOrm().LoadRelated(v, "Images")
			if len(v.Images) > 0 {
				tempdata.Image = v.Images[0].Url
			} else {
				tempdata.Image = "/static/images/bitimage.png"
			}
			arrBbsData11 = append(arrBbsData11, tempdata)
		}

		alldata = append(alldata, Item{
			Name: "美食",
			Item: arrBbsData11,
		})
		// 旅行
		var arrBbsData10 []Data
		var bbs10 []*models.Bbs
		qs = o.QueryTable(new(models.Bbs))
		qs = qs.Filter("status", 1)
		qs = qs.Filter("category_id__in", 10)
		qs = qs.Filter("Customer__Username__isnull", false)
		_, err = qs.OrderBy("-id").RelatedSel().Limit(10).All(&bbs10)
		for _, v := range bbs10 {
			var tempdata Data
			tempdata.Title = v.Content
			tempdata.Url = "/bbs/detail/" + fmt.Sprintf("%d", v.Id) + ".html"
			tempdata.Date = v.Created
			tempdata.Num = v.Review
			orm.NewOrm().LoadRelated(v, "Images")
			if len(v.Images) > 0 {
				tempdata.Image = v.Images[0].Url
			} else {
				tempdata.Image = "/static/images/bitimage.png"
			}
			arrBbsData10 = append(arrBbsData10, tempdata)
		}

		alldata = append(alldata, Item{
			Name: "旅行",
			Item: arrBbsData10,
		})
		// 读书
		var arrBbsData9 []Data
		var bbs9 []*models.Bbs
		qs = o.QueryTable(new(models.Bbs))
		qs = qs.Filter("status", 1)
		qs = qs.Filter("category_id__in", 9)
		qs = qs.Filter("Customer__Username__isnull", false)
		_, err = qs.OrderBy("-id").RelatedSel().Limit(10).All(&bbs9)
		for _, v := range bbs9 {
			var tempdata Data
			tempdata.Title = v.Content
			tempdata.Url = "/bbs/detail/" + fmt.Sprintf("%d", v.Id) + ".html"
			tempdata.Date = v.Created
			tempdata.Num = v.Review
			orm.NewOrm().LoadRelated(v, "Images")
			if len(v.Images) > 0 {
				tempdata.Image = v.Images[0].Url
			} else {
				tempdata.Image = "/static/images/bitimage.png"
			}
			arrBbsData9 = append(arrBbsData9, tempdata)
		}

		alldata = append(alldata, Item{
			Name: "读书",
			Item: arrBbsData9,
		})
		// 穿搭
		var arrBbsData8 []Data
		var bbs8 []*models.Bbs
		qs = o.QueryTable(new(models.Bbs))
		qs = qs.Filter("status", 1)
		qs = qs.Filter("category_id__in", 8)
		qs = qs.Filter("Customer__Username__isnull", false)
		_, err = qs.OrderBy("-id").RelatedSel().Limit(10).All(&bbs8)
		for _, v := range bbs8 {
			var tempdata Data
			tempdata.Title = v.Content
			tempdata.Url = "/bbs/detail/" + fmt.Sprintf("%d", v.Id) + ".html"
			tempdata.Date = v.Created
			tempdata.Num = v.Review
			orm.NewOrm().LoadRelated(v, "Images")
			if len(v.Images) > 0 {
				tempdata.Image = v.Images[0].Url
			} else {
				tempdata.Image = "/static/images/bitimage.png"
			}
			arrBbsData8 = append(arrBbsData8, tempdata)
		}

		alldata = append(alldata, Item{
			Name: "穿搭",
			Item: arrBbsData8,
		})
		// 宠物
		var arrBbsData7 []Data
		var bbs7 []*models.Bbs
		qs = o.QueryTable(new(models.Bbs))
		qs = qs.Filter("status", 1)
		qs = qs.Filter("category_id__in", 7)
		qs = qs.Filter("Customer__Username__isnull", false)
		_, err = qs.OrderBy("-id").RelatedSel().Limit(10).All(&bbs7)
		for _, v := range bbs7 {
			var tempdata Data
			tempdata.Title = v.Content
			tempdata.Url = "/bbs/detail/" + fmt.Sprintf("%d", v.Id) + ".html"
			tempdata.Date = v.Created
			tempdata.Num = v.Review
			orm.NewOrm().LoadRelated(v, "Images")
			if len(v.Images) > 0 {
				tempdata.Image = v.Images[0].Url
			} else {
				tempdata.Image = "/static/images/bitimage.png"
			}
			arrBbsData7 = append(arrBbsData7, tempdata)
		}

		alldata = append(alldata, Item{
			Name: "宠物",
			Item: arrBbsData7,
		})

		c.Data["AllData"] = alldata

		var images []*models.BbsImages
		o.QueryTable(new(models.BbsImages)).OrderBy("-id").RelatedSel().Limit(6).All(&images)
		for k, v := range images {
			c.Data["Images"+fmt.Sprintf("%d", k)] = v
		}

		todayArticleQs := o.QueryTable(new(admin.Article))
		todayArticleQs = todayArticleQs.Filter("status", 1)
		todayArticleQs = todayArticleQs.Filter("Customer__Username__isnull", false)
		todayArticleQs = todayArticleQs.Filter("Category__Name__isnull", false)
		allCountArticle, _ := todayArticleQs.Count()
		c.Data["AllCountArticle"] = allCountArticle
		todayArticleQs = todayArticleQs.Filter("created__gte", time.Now().Format("2006-01-02 00:00:00"))
		tdCountArticle, _ := todayArticleQs.Count()
		c.Data["TdCountArticle"] = tdCountArticle
		todayBbsQs := o.QueryTable(new(models.Bbs))
		todayBbsQs = todayBbsQs.Filter("status", 1)
		allCountBbs, _ := todayBbsQs.Count()
		c.Data["AllCountBbs"] = allCountBbs
		todayBbsQs = todayBbsQs.Filter("created__gte", time.Now().Format("2006-01-02 00:00:00"))
		tdCountBbs, _ := todayBbsQs.Count()
		c.Data["TdCountBbs"] = tdCountBbs
	}

	if page == "old" {
		c.TplName = "home/" + c.Template + "/oldindex.html"
	} else {

		c.TplName = "home/" + c.Template + "/index.html"
	}
}
