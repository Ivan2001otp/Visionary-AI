package scrap

import (
	"fmt"
	"log"
	"strconv"
	"time"

	model "github.com/Ivan2001otp/Visionary-AI/Model"
	util "github.com/Ivan2001otp/Visionary-AI/Util"
	"github.com/gocolly/colly"
)

func LapTopScraper() ([]model.Product, error) {
	util.Print("start..")
	var laptopList []model.Product

	const TARGET_LINK = "https://www.amazon.in/s?i=computers&rh=n%3A7198569031&fs=true&page=2&ref=sr_pg_"
	var anyErr error

	c := colly.NewCollector(
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 2 * time.Second,
	})

	c.OnRequest(func(request *colly.Request) {
		util.Print("requesting...")
	})

	c.OnResponse(func(response *colly.Response) {
		util.Print("response recieving..")
	})

	c.OnHTML("div.s-widget-container", func(jrColly *colly.HTMLElement) {

		imgUrl := jrColly.ChildAttr("img.s-image", "src")
		pname := jrColly.ChildText("a.a-link-normal > span.a-size-base-plus")
		rating_ := jrColly.ChildText("i.a-icon-star-small > span.a-icon-alt")
		gRating_ := jrColly.ChildText("a.a-link-normal > span.a-size-base")
		price_ := jrColly.ChildText("span.a-price > span.a-offscreen")
		detailPageUrl := jrColly.ChildAttr("a.a-link-normal", "href")
		// product_type := "smartwatch"
		// var i_count int = 0

		util.Print("image ->" + imgUrl)
		// util.Print("name -> " + pname)
		// util.Print(rating_)

		// util.Print(gRating_)
		// util.Print("Price:" + priceStringTrimmer(price_))
		// util.Print("Href-> " + detailPageUrl)

		if !util.IsEmpty(imgUrl) && !util.IsEmpty(pname) && !util.IsEmpty(rating_) && !util.IsEmpty(gRating_) && !util.IsEmpty(price_) && !util.IsEmpty(detailPageUrl) {

			var laptop model.Product
			laptop.ProductName = pname
			laptop.GlobalRating = gRating_
			laptop.ProductRating = rating_
			laptop.ProductDetailUrl = detailPageUrl
			laptop.ProductPrice = price_
			laptop.ProductImg = imgUrl

			laptopList = append(laptopList, laptop)
		}

	})

	maxTries := 3
	retryCount := 0
	c.OnError(func(r *colly.Response, err error) {
		anyErr = err
		util.Print("Something error!")

		//for the error -> ServiceUnavailable,blocked by ErrRobotsTxtblocked.
		if colly.ErrRobotsTxtBlocked == err {
			util.Print("Blocked by Robots.txt file -> " + err.Error())

			log.Fatal(err)

		} else if err.Error() == "Service Unavailable" {
			fmt.Println("Error while smartphone scraping->", err)

			//retry logic
			retryCount++
			if retryCount <= maxTries {
				fmt.Println("retrying ... ", retryCount)
				delay := time.Second * 2 * time.Duration(time.Duration(retryCount).Seconds())
				time.Sleep(delay)
				err := c.Visit(TARGET_LINK + strconv.Itoa(4))
				if err != nil {
					util.Print("Error captured while revisting after retry..->" + string(err.Error()))
				}
			}

		}
	})

	// var wg sync.WaitGroup

	// wg.Add(1)
	// go func() {
	// defer wg.Done()
	/*
		50 -> 843i -> 2
		1000 -> x
	*/
	for i := 1; i <= 2; i++ {
		c.Visit(TARGET_LINK + strconv.Itoa(i))
	}
	// }()

	// wg.Wait()
	c.Wait()

	if anyErr != nil {
		return nil, anyErr
	}

	if len(laptopList) == 0 {
		util.Print("not empty")
	} else {
		util.Print("empty")
	}

	return laptopList, nil
}
