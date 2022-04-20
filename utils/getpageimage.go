
package utils

import (
	"context"
	"io/ioutil"
	"log"
	// "time"
	"github.com/chromedp/chromedp"
	// "github.com/chromedp/chromedp/device"
)

func GetPageImage(path,url string) string {
	// create context

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	go func() {
		defer func() {
		  if e := recover(); e != nil {
			log.Printf("recover: %v\n", e)
		  }
		}()
	}()
	
	// ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
	// defer cancel()
	// run
	var b1 []byte
	if err := chromedp.Run(ctx,
		//chromedp.Emulate(device.BlackBerryZ30landscape),
		chromedp.EmulateViewport(1560, 877),
		chromedp.Navigate(url),
		// chromedp.Sleep(2 * time.Second),
		chromedp.CaptureScreenshot(&b1),
	);err != nil {
		log.Fatal(err)
		return ""
	}

	if err := ioutil.WriteFile(path, b1, 0o644); err != nil {
		log.Fatal(err)
		return ""
	}

	//log.Printf("wrote screenshot1.png and screenshot2.png")
	return path
}
