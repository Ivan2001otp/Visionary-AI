package scrap

import (
	"fmt"

	util "github.com/Ivan2001otp/Visionary-AI/Util"
	"github.com/gocolly/colly"
)

type LapTop struct {
	ProductName      string
	ProductType      string
	ProductImg       string
	ProductDetailUrl string
	ProductRating    string
	GlobalRating     string
	ProductPrice     string
	ProductRetailers string
}

func LapTopScraper() ([]LapTop, error) {
	util.Print("start..")
	var laptopList []LapTop

	const TARGET_LINK = "https://www.amazon.in/s?i=computers&rh=n%3A7198569031&fs=true&page=2&ref=sr_pg_1"
	var anyErr error

	c := colly.NewCollector()

	c.Visit(TARGET_LINK)

	c.OnRequest(func(request *colly.Request) {
		util.Print("requesting...")
	})

	c.OnResponse(func(response *colly.Response) {
		util.Print("response recieving..")
	})

	c.OnHTML("div.s-card-container", func(jrColly *colly.HTMLElement) {
		util.Print("Started to scrap...")
		var laptop LapTop

		imgUrl := jrColly.ChildAttr("img.s-image", "src")
		pname := jrColly.ChildText("a.a-link-normal > span.a-text-normal")
		rating_ := jrColly.ChildText("i.a-icon-star-small > span.a-icon-alt")
		gRating_ := jrColly.ChildText("a.a-link-normal > span.a-size-base")
		price_ := jrColly.ChildText("span.a-price > span.a-offscreen")
		detailPageUrl := jrColly.ChildAttr("a.a-link-normal", "href")
		product_type := "smartwatch"
		var i_count int = 0

		util.Print(imgUrl)
		util.Print(pname)
		util.Print(rating_)
		util.Print(gRating_)
		util.Print(price_)
		util.Print(detailPageUrl)
		fmt.Println()
		fmt.Println()

		laptop.ProductType = product_type
		if !util.IsEmpty(imgUrl) {
			laptop.ProductImg = imgUrl
			i_count++
		}

		if !util.IsEmpty(pname) {
			laptop.ProductName = pname
			i_count++
		}
		if !util.IsEmpty(rating_) {
			laptop.ProductRating = rating_
			i_count++
		}

		if !util.IsEmpty(gRating_) {
			laptop.GlobalRating = gRating_
			i_count++
		}

		if !util.IsEmpty(price_) {
			laptop.ProductPrice = price_
			i_count++
		}

		if !util.IsEmpty(detailPageUrl) {
			laptop.ProductDetailUrl = detailPageUrl
			i_count++
		}

		// if i_count >= 6 {
		laptopList = append(laptopList, laptop)
		// }
		util.Print("Don with scrap...")

	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error while smartphone scraping->", err)
		anyErr = err
	})

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
