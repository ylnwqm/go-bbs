package utils
import (

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/ulule/limiter/v3"
	"net/http"
	"strconv"
	"strings"
)

type RateLimiter struct {
	GeneralLimiter *limiter.Limiter
	LoginLimiter   *limiter.Limiter
}


func RateLimit(r *RateLimiter, ctx *context.Context) {
	var (
		limiterCtx limiter.Context
		ip         string
		err        error
		req        = ctx.Request
	)



	//ipVerifyList := "66.249.71.0-66.249.71.255"
	// blacklistip := ctx.Input.IP()
	// ipSlice := strings.Split(ipVerifyList, `-`)
	// if (ip2Int(blacklistip) >= ip2Int(ipSlice[0]) && ip2Int(blacklistip) <= ip2Int(ipSlice[1])) || blacklistip == "194.209.25.132"{
	// 	logs.Debug("你的IP " + blacklistip + " 存在可疑行为已被拉入黑名单，如有误报请联系站长(1920853199@qq.com)")
	// 	ctx.Abort(http.StatusForbidden, "403")
	// 	return
	// }

	//fmt.Printf("%s\n",blacklistip)

	if strings.HasPrefix(ctx.Input.URL(), "/login") {
		ip = ctx.Input.IP()
		limiterCtx, err = r.LoginLimiter.Get(req.Context(), ip)
	} else {
		ip = ctx.Input.IP()
		limiterCtx, err = r.GeneralLimiter.Get(req.Context(), ip)
	}
	if err != nil {
		ctx.Abort(http.StatusInternalServerError, err.Error())
		return
	}

	h := ctx.ResponseWriter.Header()
	h.Add("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
	h.Add("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
	h.Add("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

	if limiterCtx.Reached {
		logs.Debug("Too Many Requests from %s on %s", ip, ctx.Input.URL())
		//refer to https://beego.me/docs/mvc/controller/errors.md for error handling
		ctx.Abort(http.StatusTooManyRequests, "429")
		return
	}

}

func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func ip2Int(ip string) int64 {
	if len(ip) == 0 {
		return 0
	}
	bits := strings.Split(ip, ".")
	if len(bits) < 4 {
		return 0
	}
	b0 := string2Int(bits[0])
	b1 := string2Int(bits[1])
	b2 := string2Int(bits[2])
	b3 := string2Int(bits[3])

	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func string2Int(in string) (out int) {
	out, _ = strconv.Atoi(in)
	return
}