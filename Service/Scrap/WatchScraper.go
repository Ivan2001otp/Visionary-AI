package scrap

import (
	"fmt"
	"strconv"

	util "github.com/Ivan2001otp/Visionary-AI/Util"
	"github.com/gocolly/colly"
)

func WatchScraper() (string, error) {

	c := colly.NewCollector(
		// colly.Async(true),
		colly.UserAgent(util.RandomUserAgent()),
	)

	// c.Limit(&colly.LimitRule{
	// 	Delay:       1 * time.Second,
	// 	RandomDelay: 500 * time.Millisecond,
	// })

	var anyErr error

	c.OnError(func(r *colly.Response, err error) {
		if err != nil {
			util.Print(err.Error())
			anyErr = err
		}
	})

	c.OnHTML("div.s-card-container", func(jrColly *colly.HTMLElement) {
		fmt.Println("exe")
		imgUrl := jrColly.ChildAttr("img.s-image", "src")

		fmt.Println("Img -> ", imgUrl)
	})

	for i := 1; i <= 3; i++ {
		c.Visit("https://www.amazon.in/s?k=kitchen+groceries&i=kitchen&rh=n%3A1379989031&ref=sr_pg_" + strconv.Itoa(i))
	}

	if anyErr != nil {
		return "", anyErr
	}

	return "success", nil
}
