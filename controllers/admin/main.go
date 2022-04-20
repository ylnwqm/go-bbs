package admin

import (
	"go-bbs/utils"
	"go-bbs/utils/sys"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type Log struct {
	Num  int
	Date string
	Ip   []string
}

type MainController struct {
	BaseController
}

func (c *MainController) Index() {
	c.TplName = "admin/index.html"
}

func (c *MainController) Welcome() {

	o := orm.NewOrm()

	var log []orm.Params

	// time.Now().AddDate(0,0,-7) 七天前
	o.Raw("SELECT `id`,`ip`,`create` FROM `log` where `create` >= ? and `ip` <> '152.32.212.226' ORDER BY `create` ASC", time.Now().AddDate(0, 0, -30)).Values(&log)

	var pv = make(map[string]Log)
	var uv = make(map[string]Log)
	var dateSlice []string
	var pvSlice []int
	var uvSlice []int

	var keys []string

	for _, v := range log {
		//utils.StringToTime(v["create"])
		// 获取日期
		var key = utils.StringToTime(v["create"]).Format("2006-01-02")

		var flag = true
		for _, k := range keys {
			if k == key {
				flag = false
			}
		}
		if flag {
			keys = append(keys, key)
		}

		// 统计pv
		pvTemp := pv[key]
		pv[key] = Log{
			Num: pvTemp.Num + 1,
		}
		// 统计uv
		uvTemp := uv[key]
		var uvFlag = true
		for _, k := range uvTemp.Ip {
			if k == v["ip"] {
				uvFlag = false
			}
		}
		if uvFlag {
			uv[key] = Log{
				Num: uvTemp.Num + 1,
				Ip:  append(uvTemp.Ip, v["ip"].(string)),
			}
		}

	}

	for _, k := range keys {
		dateSlice = append(dateSlice, k)
		pvSlice = append(pvSlice, pv[k].Num)
		uvSlice = append(uvSlice, uv[k].Num)
	}

	c.Data["Date"] = dateSlice
	c.Data["Pv"] = pvSlice
	c.Data["Uv"] = uvSlice

	df, _ := sys.Df()
	c.Data["Df"] = df

	var customer []orm.Params

	// time.Now().AddDate(0,0,-7) 七天前
	// SELECT date_format(`created`, '%Y-%m-%d') as t, count(id) as c FROM `customer` group by date_format(`created`, '%Y-%m-%d') ORDER BY `id` ASC;
	o.Raw("SELECT date_format(`created`, '%Y-%m-%d') as t, count(id) as c FROM `customer` where `created` >= '1970-01-01 00:00:00' and `email` != 'null' group by date_format(`created`, '%Y-%m-%d') ORDER BY `id` ASC").Values(&customer)

	var userNum = make([]int, 0, len(customer))
	var userTime = make([]string, 0, len(customer))
	var countUserNum int
	for _, v := range customer {
		userTime = append(userTime, v["t"].(string))
		userNum_int, _ := strconv.Atoi(v["c"].(string))
		countUserNum += userNum_int
		userNum = append(userNum, countUserNum)
	}

	c.Data["UserNum"] = userNum
	c.Data["UserTime"] = userTime

	var referer []orm.Params

	o.Raw(`SELECT
    COUNT(id) as c,
         CASE 
         WHEN referer RLIKE "^http://|^https://" THEN SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING_INDEX(referer, "/", 3),"/",-1),".",-2)
         ELSE SUBSTRING_INDEX(referer, "/", 1)
		 END as referer
FROM
    log
GROUP BY
         CASE 
         WHEN referer RLIKE "^http://|^https://" THEN SUBSTRING_INDEX(SUBSTRING_INDEX(SUBSTRING_INDEX(referer, "/", 3),"/",-1),".",-2)
         ELSE SUBSTRING_INDEX(referer, "/", 1)
		 END`).Values(&referer)

	var rTitle = make([]string, 1)
	//var rv = make(map[string]interface{})
	var item = make([]map[string]interface{}, 1)
	for _, v := range referer {
		if v["referer"].(string) == "" {
			continue
			v["referer"] = "直接输入网址或书签"
		}
		if v["referer"].(string) == "nihongdengxia.com" {
			continue
		}

		c, err := strconv.Atoi(v["c"].(string))
		if err != nil {
			continue
		}

		if c < 100 {
			continue
		}

		rTitle = append(rTitle, v["referer"].(string))

		item = append(item, map[string]interface{}{
			"value": v["c"],
			"name":  v["referer"],
		})
	}

	c.Data["RTitle"] = rTitle
	c.Data["Referer"] = item

	c.TplName = "admin/welcome.html"
}
