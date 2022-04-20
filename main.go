package main

import (
	"flag"
	_ "go-bbs/routers"
	"go-bbs/utils"
	"go-bbs/utils/syncdb"
	"os"

	// "go-bbs/utils/sitemap"
	"go-bbs/utils/hotnews"
	"html/template"
	"net/http"
	_ "net/http/pprof"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/robfig/cron"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	// "github.com/garyburd/redigo/redis"
	// "fmt"
	// "time"
)

func init() {
	var t bool
	var install bool
	flag.BoolVar(&t, "t", false, "爬虫定时任务是否开启")
	flag.BoolVar(&install, "install", false, "是否安装")
	flag.Parse()
	args := os.Args
	for _, v := range args {
		if v == "--install" {
			syncdb.Syncdb()
			os.Exit(0)
		}
	}

	syncdb.Connect()

	beego.AddFuncMap("IndexForOne", utils.IndexForOne)
	beego.AddFuncMap("IndexAddOne", utils.IndexAddOne)
	beego.AddFuncMap("IndexDecrOne", utils.IndexDecrOne)
	beego.AddFuncMap("StringReplace", utils.StringReplace)
	beego.AddFuncMap("TimeStampToTime", utils.TimeStampToTime)
	beego.AddFuncMap("TemlpateTime", utils.TemlpateTime)
	beego.AddFuncMap("IsOdd", func(num int) string {
		if num%2 == 0 {
			return "odd"
		}
		return "even"
	})

	if t == true {
		//每天0点定时更新站点地图
		go func() {
			c := cron.New()
			//*/1 0 * * *
			// 0 0 * * *
			// c.AddFunc("@daily", func() {
			// 	sitemap.Sitemap("./", conf.String("url"))
			// })

			c.AddFunc("0 */25 * * *", func() {
				hp := hotnews.Hupu{}
				hp.Get("https://www.hupu.com/")
			})

			c.AddFunc("0 */15 * * *", func() {
				w := hotnews.Weibo{}
				w.Get("https://s.weibo.com/top/summary")
			})

			c.AddFunc("0 */30 * * *", func() {
				b := hotnews.Baidu{}
				b.Get("https://top.baidu.com/board?tab=realtime")
			})

			c.AddFunc("0 */60 * * *", func() {
				k := hotnews.Kaolamedia{}
				k.Get("https://www.kaolamedia.com/hot")
			})

			c.Start()
		}()
	}

	// debug sql
	// orm.Debug = true

}

func page_not_found(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("404.html").ParseFiles("views/home/nihongdengxia/404.html")
	data := make(map[string]interface{})
	data["content"] = template.HTML("404 Page Not Found")
	data["url"] = "/"
	t.Execute(rw, data)

}

func main() {
	// mount /dev/vdb /work
	// docker build -t chromegin . && docker run -p 6666:6666 -v /work/chromegin/chrome_screen_shot:/data --name chromegin chromegin
	r := &utils.RateLimiter{}

	rate, err := limiter.NewRateFromFormatted("5000-D")

	utils.PanicOnError(err)
	r.GeneralLimiter = limiter.New(memory.NewStore(), rate)

	loginRate, err := limiter.NewRateFromFormatted("10-M")
	utils.PanicOnError(err)
	r.LoginLimiter = limiter.New(memory.NewStore(), loginRate)

	//More on Beego filters here https://beego.me/docs/mvc/controller/filter.md
	beego.InsertFilter("/*", beego.BeforeRouter, func(c *context.Context) {
		utils.RateLimit(r, c)
	}, true)

	//refer to https://beego.me/docs/mvc/controller/errors.md for error handling
	beego.ErrorHandler("429", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too Many Requests"))
		return
	})

	// 黑名单
	beego.ErrorHandler("403", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("你的IP存在可疑行为已被拉入黑名单，如有误报请联系站长(1920853199@qq.com)"))
		return
	})

	beego.ErrorHandler("404", page_not_found)
	//bee generate appcode -tables="cron" -driver=mysql -conn="root:root@tcp(127.0.0.1:3306)/blog" -level=3
	beego.Run()

	// s := utils.MusicGet(2175288)
	// fmt.Printf("%v", s)
	// path, err := utils.DownImage("https://img9.doubanio.com/view/subject/s/public/s3745215.jpg")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(path)
	// fmt.Println(utils.TemlpateTime(time.Now()))
	// fmt.Println(utils.TemlpateTime(time.Now().Add(-time.Second * 40)))
	// fmt.Println(utils.TemlpateTime(time.Now().Add(-time.Second * 10)))
	// fmt.Println(utils.TemlpateTime(time.Now().Add(-time.Second * 60 * 2)))
	// fmt.Println(utils.TemlpateTime(time.Now().Add(-time.Hour * 2)))
	// fmt.Println(utils.TemlpateTime(time.Now().Add(-time.Hour * 26)))
	// fmt.Println(utils.TemlpateTime(time.Now().Add(-time.Hour * 79)))
	// fmt.Println(utils.TemlpateTime(time.Now().Add(-time.Hour * 24 * 31 * 3)))
	// fmt.Println(utils.TemlpateTime(time.Now().Add(-time.Hour * 24 * 365 * 2)))
	// fmt.Println(utils.TemlpateTime(time.Now().Add(-time.Hour * 24 * 365 * 200)))
	// bytes, err := utils.GenDefaultQRCode("Hello World")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// err = ioutil.WriteFile("1.jpg", bytes, 0666)

	// bytes, err = utils.GenDefaultQRCode("http://www.baidu.com")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// err = ioutil.WriteFile("2.jpg", bytes, 0666)

	//worker, _ := utils.NewIdWorker(1)
	//ID, _ := worker.GetNextId()
	// username := utils.CreateUsername()
	//fmt.Println(ID)
	// ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("p-id", "P66-"))
	// id := utils.NextID(ctx)
	// fmt.Println(id)
	// fmt.Println(utils.NextID(context.Background()))
	// fmt.Println(utils.NextNumID())
	// fmt.Println(utils.UUID())
	//fmt.Println(utils.UUIDShort())
	// utils.SendEmail(utils.Email{
	// 	From: "1920853199@qq.com",
	// 	To:   "491126240@qq.com",
	// 	Header: map[string]string{
	// 		"Subject": "霓虹灯下注册验证码",
	// 	},
	// 	Template: "views/home/nihongdengxia/reg-email.html",
	// 	Data: utils.Reg{
	// 		Tag:  "注册",
	// 		Code: "123",
	// 		Date: time.Now().Format("2006-01-02 15:04:05"),
	// 	},
	// })

}
