package scrap

import (
	"fmt"
	"strconv"
	"time"

	model "github.com/Ivan2001otp/Visionary-AI/Model"
	retailer "github.com/Ivan2001otp/Visionary-AI/Retailer"
	util "github.com/Ivan2001otp/Visionary-AI/Util"

	"github.com/gocolly/colly"
)

/*
Note : Do not alter the LINK variables... otherwise entire code will break;
*/
func WatchScraper() ([]model.Product, error) {
	var watchList []model.Product
	var prevPageHolder int = 1

	const LINK string = "https://www.amazon.in/s?k=smart+watch&i=electronics&rh=n%3A976419031%2Cp_123%3A214020%7C230542%7C358517%7C370584%7C42717%7C851753&dc&rnid=91049095031&ref=sr_pg_"
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent(util.RandomUserAgent()),
	)

	c.Limit(&colly.LimitRule{
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	var anyErr error

	maxTry := 3
	retryCount := 0
	c.OnError(func(r *colly.Response, err error) {
		if err != nil {
			util.Print(err.Error())
			anyErr = err

			//retry logic for 3 times..
			if err.Error() == "Service Unavailable" {
				retryCount++
				if retryCount <= maxTry {
					fmt.Println("Retrying the cycle ", retryCount)

					//later we can apply here for loop to scrap from prevpage to next following pages.
					for i := prevPageHolder + 1; i <= prevPageHolder+2; i++ {

						c.Visit(LINK + strconv.Itoa(i))

					}
				}
			} else if err == colly.ErrRobotsTxtBlocked {
				fmt.Println("Violated Robots.Txt file...")
			}
		}
	})

	c.OnHTML("div.s-card-container", func(jrColly *colly.HTMLElement) {
		var item model.Product

		imgUrl := jrColly.ChildAttr("img.s-image", "src")
		hrefLink := jrColly.ChildAttr("a.a-link-normal", "href")
		_rating := jrColly.ChildText("i.a-icon-star-small > span.a-icon-alt")
		_gRating := jrColly.ChildText("a.a-link-normal > span.a-size-base")
		_price := jrColly.ChildText("span.a-price > span.a-offscreen")
		_name := jrColly.ChildText("a.a-link-normal > span.a-size-base-plus")

		if !util.IsEmpty(imgUrl) && !util.IsEmpty(hrefLink) && !util.IsEmpty(_rating) && !util.IsEmpty(_gRating) && !util.IsEmpty(_price) && !util.IsEmpty(_name) {
			// util.Print("not empty")
			item.ProductImg = imgUrl
			// util.Print("img from item-> " + item.ProductImg)

			item.ProductDetailUrl = hrefLink
			item.ProductRating = _rating
			item.GlobalRating = _gRating
			item.ProductPrice = util.PriceStringTrimmer(_price)
			item.ProductName = _name
			item.CategoryType = "smartwatch"
			item.ProductRetailers = retailer.WatchRetailer(_name)

			watchList = append(watchList, item)
		}
		// fmt.Println("Img for real -> ", imgUrl)
		// fmt.Println("HrefLink -> ", item.ProductDetailUrl)
		// fmt.Println("rating -> ", item.ProductRating)
		// fmt.Println("gRating -> ", item.GlobalRating)
		// fmt.Println("price -> ", item.ProductPrice)
		fmt.Println("name -> ", item.ProductName)
		fmt.Println()
	})

	for i := 1; i <= 2; i++ {
		c.Visit(LINK + strconv.Itoa(i))
		prevPageHolder = i
	}
	c.Wait()

	if anyErr != nil {
		return nil, anyErr
	}

	return watchList, nil
}
