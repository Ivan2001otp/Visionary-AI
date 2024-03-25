package scrap

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	model "github.com/Ivan2001otp/Visionary-AI/Model"
	retailer "github.com/Ivan2001otp/Visionary-AI/Retailer"
	util "github.com/Ivan2001otp/Visionary-AI/Util"
	"github.com/gocolly/colly"
)

func TelevisionScraper() ([]model.Product, error) {
	const LINK1 string = "https://www.amazon.in/s?i=electronics&bbn=976419031&rh=n%3A12045104031&ref=sr_pg_"
	var televisionList []model.Product

	var anyErr error
	var prevPageHolder int = 1

	c := colly.NewCollector(
		colly.UserAgent(util.RandomUserAgent()),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	/*
		c.OnRequest(func(r *colly.Request) {})
		c.OnResponse(func(r *colly.Response) {})
		c.OnScraped(func(r *colly.Response) {})
	*/
	c.OnHTML("div.s-result-item", func(h *colly.HTMLElement) {

		// p := h.Response
		// i += 1
		// fmt.Println(i, "->", p.StatusCode)
		var productClone model.Product

		imgUrl := h.ChildAttr("img.s-image", "src")                        //img
		hrefLink := h.ChildAttr("a.a-link-normal", "href")                 //href
		pname := h.ChildText("a.a-link-normal > span.a-text-normal")       // pname
		rating := h.ChildText("i.a-icon-star-small > span.a-icon-alt")     //rating out 5.
		price := h.ChildText("span.a-price > span.a-offscreen")            //price.
		globalRating_ := h.ChildText("a.a-link-normal > span.a-size-base") //global rating

		fmt.Println("Television name :", pname)

		if !util.IsEmpty(imgUrl) && !util.IsEmpty(hrefLink) && !util.IsEmpty(pname) && !util.IsEmpty(rating) && !util.IsEmpty(price) && !util.IsEmpty(globalRating_) {

			if !strings.HasPrefix(hrefLink, "https://www.amazon.in") {
				hrefLink = "https://www.amazon.in" + hrefLink
			}
			productClone.ProductDetailUrl = hrefLink

			productClone.ProductRetailers = retailer.TelevisionRetailers()
			productClone.CategoryType = "television"
			productClone.ProductImg = imgUrl
			productClone.ProductName = pname
			productClone.ProductRating = rating
			productClone.ProductPrice = price
			productClone.GlobalRating = globalRating_

			televisionList = append(televisionList, productClone)
		}

	})

	maxTry := 3
	retryCount := 0
	c.OnError(func(r *colly.Response, err error) {

		if err != nil {
			anyErr = err
			fmt.Println("Error catched while scraping television data..", err)

			if err.Error() == "Service Unavailable" {
				retryCount++
				if retryCount <= maxTry {
					fmt.Println("Retrying to scrap television data..", retryCount)
					//tries to scrap few pages i.e, 3 pages more.
					for i := prevPageHolder + 1; i <= prevPageHolder+2; i++ {
						c.Visit(LINK1 + strconv.Itoa(i))
					}
				}
			} else if err.Error() == colly.ErrRobotsTxtBlocked.Error() {
				fmt.Println("Violates Robots.txt file!")
			}
		}
	})

	for i := 1; i <= 2; i++ {
		c.Visit(LINK1 + strconv.Itoa(i))
		prevPageHolder = i
	}

	if anyErr != nil {
		return nil, anyErr
	}

	return televisionList, nil

}
