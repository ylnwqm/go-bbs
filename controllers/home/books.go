package home

import (
	"encoding/json"
	"fmt"
	"go-bbs/models"
	"go-bbs/models/admin"
	"go-bbs/utils"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type BooksController struct {
	BaseController
}

type Rating struct {
	Max       string `json:"max"`
	NumRaters string `json:"numRaters"`
	Average   string `json:"average"`
	Min       string `json:"min"`
}
type Books struct {
	EbookUrl    string   `json:"ebook_url"`
	Author      []string `json:"author"`
	Pubdate     string   `json:"pubdate"`
	Catalog     string   `json:"catalog"`
	Pages       string   `json:"pages"`
	Title       string   `json:"title"`
	AuthorIntro string   `json:"author_intro"`
	Summary     string   `json:"summary"`
	Rating      Rating   `json:"rating"`
	Id          string   `json:"id"`
	Isbn10      string   `json:"isbn10"`
	Isbn13      string   `json:"isbn13"`
	Image       string   `json:"image"`
}
type HttpResponse struct {
	Count int     `json:"count"`
	Start int     `json:"start"`
	Total int     `json:"total"`
	Books []Books `json:"books"`
}

var ReadSatus = map[int]string{
	1: "正在阅读",
	2: "准备阅读",
	3: "已经读完",
	0: "删除",
}

func (c *BooksController) Books() {
	year := c.GetString("year", time.Now().Format("2006"))

	yearStr := fmt.Sprintf("%s年", year)
	//fmt.Println(year)
	yearIntTmp, _ := strconv.Atoi(year)
	var oldYear bool
	if yearIntTmp <= 2021 {
		oldYear = true
		yearStr = "2022年之前"
	}
	nextYear := fmt.Sprintf("%d", yearIntTmp+1)

	customer_id, _ := c.GetInt("customer_id")
	if customer_id == 0 {
		customer_id = 653
	}

	type book struct {
		SkuId      int
		GoodsId    int
		Title      string
		Image      string
		Url        string
		AttrString string
		Attr       map[string]string
		EbookUrl   string
		Id         int
	}
	o := orm.NewOrm()
	// 查询所有的年份
	var allYearbooks []*models.Books
	var allYear = make(map[string]string)
	o.QueryTable(new(models.Books)).Filter("customer_id", customer_id).Filter("status__gt", 0).All(&allYearbooks, "created")
	for _, tmpBooks := range allYearbooks {
		tyear := tmpBooks.Created.Format("2006")
		tyearIntTmp, _ := strconv.Atoi(tyear)
		if tyearIntTmp <= 2021 {
			allYear["2022年之前"] = "2021"
		} else {
			allYear[fmt.Sprintf("%d", tyearIntTmp)] = fmt.Sprintf("%d", tyearIntTmp)
		}
	}

	// 统计书籍
	var allbooks []*models.Books
	allCount, _ := o.QueryTable(new(models.Books)).Filter("customer_id", customer_id).Filter("status", 3).Count()

	// 查询当年书籍
	qs := o.QueryTable(new(models.Books)).Filter("customer_id", customer_id).Filter("status__gt", 0)

	if oldYear {
		qs = qs.Filter("Created__lt", fmt.Sprintf("%s-01-01 00:00:00", nextYear)).OrderBy("-Created")
	} else {
		qs = qs.Filter("Created__gte", fmt.Sprintf("%s-01-01 00:00:00", year)).Filter("Created__lt", fmt.Sprintf("%s-01-01 00:00:00", nextYear)).OrderBy("-Created")
	}
	// yearCount, _ := qs.Count()
	qs.All(&allbooks)

	var responseData = map[string]map[string][]book{
		"已经读完": make(map[string][]book),
		"准备阅读": make(map[string][]book),
		"正在阅读": make(map[string][]book),
	}

	var sortDateMap = make(map[string]bool)
	var sortDateSclie []string
	var yearReadedCount int
	var yearReadingCount int
	for _, statusBooks := range allbooks {
		if statusBooks.Status == 0 {
			continue
		}

		if statusBooks.Status == 3 {
			yearReadedCount++
		}

		if statusBooks.Status == 1 {
			yearReadingCount++
		}
		dateformat := statusBooks.Created.Format("2006/01")
		sortDateMap[dateformat] = true
		responseData[ReadSatus[statusBooks.Status]][dateformat] = append(responseData[ReadSatus[statusBooks.Status]][dateformat], book{
			Title:    statusBooks.Title,
			Image:    statusBooks.Image,
			EbookUrl: statusBooks.EbookUrl,
			Id:       statusBooks.Id,
		})
	}

	for k := range sortDateMap {
		sortDateSclie = append(sortDateSclie, k)
	}

	c.Data["Books"] = responseData
	var BooksSort []map[string]map[string][]book
	if year == time.Now().Format("2006") {
		BooksSort = []map[string]map[string][]book{
			{"正在阅读": responseData["正在阅读"]},
			{"准备阅读": responseData["准备阅读"]},
			{"已经读完": responseData["已经读完"]},
		}
	} else {
		BooksSort = []map[string]map[string][]book{
			{"已经读完": responseData["已经读完"]},
		}
	}

	// 查看统计点赞
	likeCount, _ := o.QueryTable(new(models.BooksLike)).Filter("books_customer_id", customer_id).Count()
	c.Data["LikeCount"] = likeCount

	// 是否点赞
	exist := o.QueryTable(new(models.BooksLike)).Filter("books_customer_id", customer_id).Filter("customer_id", c.CustomerId).Exist()
	c.Data["Exist"] = exist

	c.Data["Year"] = year       // 当前年
	c.Data["AllYear"] = allYear // 所有年份
	c.Data["YearStr"] = yearStr // 年份标志

	var user admin.Customer
	o.QueryTable(new(admin.Customer)).Filter("id", customer_id).One(&user, "username")
	c.Data["Customerid"] = customer_id
	c.Data["Username"] = user.Username

	c.Data["YearReadedCount"] = yearReadedCount
	c.Data["YearReadingCount"] = yearReadingCount
	c.Data["AllYearCount"] = allCount
	c.Data["BooksSort"] = BooksSort
	c.Data["SortDateSclie"] = sortDateSclie
	c.TplName = "home/" + c.Template + "/books.html"
}

const DOUBANAPIURL = "https://api.douban.com/v2/book/search"
const DOUBANAPIURLAPIKEY = "0b2bdeda43b5688921839c8ecb20399b"

func (c *BooksController) SearchBooks() {

	response := make(map[string]interface{})
	q := c.GetString("q")

	if q == "" {
		response["msg"] = "非法操作！"
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}
	var res HttpResponse

	err := utils.HTTPRequest(fmt.Sprintf("%s?q=%s&apikey=%s", DOUBANAPIURL, q, DOUBANAPIURLAPIKEY), "POST", []byte{}, &utils.RequestParams{
		Timeout: 60,
	}, &res, false)

	if err != nil {
		response["msg"] = err.Error()
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	for k, v := range res.Books {

		if v.Image != "" {
			img, err := utils.DownImage(v.Image)
			if err == nil {
				res.Books[k].Image = img
			}
		}

	}

	response["Data"] = res.Books
	response["msg"] = "Success."
	response["code"] = 200
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *BooksController) SelectBooks() {
	c.TplName = "home/" + c.Template + "/selectbooks.html"
}

func (c *BooksController) SaveBooks() {

	response := make(map[string]interface{})

	if !c.IsLogin {
		response["msg"] = "请先登录！"
		response["code"] = 500
		response["err"] = "请先登录！"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	if c.CustomerId == 0 {
		c.CustomerId = 653
	}
	data := c.Ctx.Input.RequestBody

	type RequestData struct {
		Status string  `json:"status"`
		Date   string  `json:"date"`
		Books  []Books `json:"books"`
	}
	//fmt.Printf("%v\n", data)
	var requestData RequestData
	err := json.Unmarshal(data, &requestData)

	//fmt.Printf("%v\n", requestData)
	if err != nil {
		response["msg"] = "数据有误"
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	if requestData.Status == "" || requestData.Date == "" || len(requestData.Books) == 0 {
		response["msg"] = "数据有误"
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	books := []models.Books{}

	status, _ := strconv.Atoi(requestData.Status)
	t, _ := time.ParseInLocation("2006/01/02", requestData.Date, time.Local)

	for _, v := range requestData.Books {
		ext, _ := json.Marshal(v)
		books = append(books, models.Books{
			EbookUrl:      v.EbookUrl,
			Image:         v.Image,
			Author:        v.Author[0],
			Pubdate:       v.Pubdate,
			Catalog:       v.Catalog,
			Pages:         v.Pages,
			Title:         v.Title,
			AuthorIntro:   v.AuthorIntro,
			Summary:       v.Summary,
			MaxRaters:     v.Rating.Max,
			NumRaters:     v.Rating.NumRaters,
			AverageRaters: v.Rating.Average,
			MinRaters:     v.Rating.Min,
			Isbn10:        v.Isbn10,
			Isbn13:        v.Isbn13,
			Ext:           string(ext),
			CustomerId:    c.CustomerId,
			Status:        status,
			Created:       t,
		})
	}

	o := orm.NewOrm()
	_, err = o.InsertMulti(len(books), books)
	if err != nil {
		response["msg"] = "操作失败，请联系管理员！" + err.Error()
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}
	//response["Data"] = res.Books
	response["msg"] = "Success."
	response["code"] = 200
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *BooksController) SetStatus() {
	response := make(map[string]interface{})

	if !c.IsLogin {
		response["msg"] = "请先登录！"
		response["code"] = 500
		response["err"] = "请先登录！"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	if c.CustomerId == 0 {
		c.CustomerId = 653
	}

	id, _ := c.GetInt("id", 0)
	status, _ := c.GetInt("status")

	if id == 0 {
		response["msg"] = "数据有误"
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	book, err := models.GetBooksById(id)
	if err != nil || book.CustomerId != c.CustomerId {
		response["msg"] = "非法操作"
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	book.Status = status
	book.Created = time.Now()
	err = models.UpdateBooksById(book)

	if err != nil {
		response["msg"] = "状态设置失败"
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	response["msg"] = "Success."
	response["code"] = 200
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *BooksController) SetLike() {

	response := make(map[string]interface{})

	if !c.IsLogin {
		response["msg"] = "请先登录！"
		response["code"] = 500
		response["err"] = "请先登录！"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	customer_id, _ := c.GetInt("customer_id")
	if customer_id == 0 {
		customer_id = 653
	}

	o := orm.NewOrm()
	exist := o.QueryTable(new(models.BooksLike)).Filter("customer_id", c.CustomerId).Filter("books_customer_id", customer_id).Exist()

	if exist {
		response["msg"] = "亲，您已经点过赞了。"
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	_, err := models.AddBooksLike(&models.BooksLike{
		CustomerId:      c.CustomerId,
		BooksCustomerId: customer_id,
		Created:         time.Now(),
	})

	if err != nil {
		response["msg"] = err.Error()
		response["code"] = 500
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
	}

	response["msg"] = "Success."
	response["code"] = 200
	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}
