package hotnews

import (
	"testing"
)

// func TestHupu(t *testing.T){
// 	hp := Hupu{}
// 	hp.Get("https://www.hupu.com/")
// }

// func TestWeibo(t *testing.T){
// 	hp := Hupu{}
// 	hp.Get("https://s.weibo.com/top/summary")
// }

// func TestBaidu(t *testing.T){
// 	b := Baidu{}
// 	b.Get("https://top.baidu.com/board?tab=realtime")
// }

// func TestKaolamedia(t *testing.T){
// 	k := Kaolamedia{}
// 	k.Get("https://www.kaolamedia.com/hot")
// }

func TestHistoday(t *testing.T){
	k := Hupu{}
	k.GetMatch("https://www.hupu.com/")
}