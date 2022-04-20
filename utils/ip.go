package utils

import (
	"net/http"
	"time"
    "encoding/json"
    "io/ioutil"
)

type Region struct {
    City string   `json:"city"`
}
func IP(ip string) string {
	
    client := http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(`http://ip-api.com/json/` + ip + `?fields=city&lang=zh-CN`)
    if err != nil {
        return ""
    }
    defer resp.Body.Close()
    html,_ := ioutil.ReadAll(resp.Body)
    var r Region
    err = json.Unmarshal(html,&r)
    if err != nil {
        return ""
    }

    return r.City
}