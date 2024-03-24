package scrap

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Ivan2001otp/Visionary-AI/Retailer"
	util "github.com/Ivan2001otp/Visionary-AI/Util"
	"github.com/gocolly/colly"
)

type SmartPhone struct {
	ProductName      string
	ProductType      string
	ProductImg       string
	ProductDetailUrl string
	ProductRating    string
	GlobalRating     string
	ProductPrice     string
	ProductRetailers string
}

/*
Note : Do not alter the LINK variables... otherwise entire code will break;
*/

func SmartPhoneScrapper() ([]SmartPhone, error) {
	fmt.Println("visiting smartphone ")
	var productList []SmartPhone
	const TARGET_LINK = "https://www.amazon.in/s?k=smartphones&i=electronics&rh=n%3A1389401031&page=2&ref=sr_pg_"
	var anyErr error

	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent(util.RandomUserAgent()),
	)

	c.Limit(&colly.LimitRule{
		Delay:       1 * time.Second,
		RandomDelay: 500 * time.Millisecond,
	})

	c.OnRequest(func(r *colly.Request) {})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("visited smartphone ")

	})
	c.OnHTML("div.s-card-container", func(jrColly *colly.HTMLElement) {
		var item SmartPhone

		imgUrl := jrColly.ChildAttr("img.s-image", "src")                  //image of smartphone
		detailLink := jrColly.ChildAttr("a.a-link-normal", "href")         //detail link
		pname := jrColly.ChildText("a.a-link-normal > span.a-text-normal") //product name
		rating := jrColly.ChildText("i.a-icon > span.a-icon-alt")          //rating.
		gRating := jrColly.ChildText("a.a-link-normal > span.a-size-base") //global rating.
		price := jrColly.ChildText("span.a-price > span.a-offscreen")      //price.

		if strings.HasPrefix(detailLink, "https://amazon.in") {
			detailLink = "https://amazon.in" + detailLink
		}

		item.GlobalRating = gRating
		item.ProductImg = imgUrl
		item.ProductDetailUrl = detailLink
		item.ProductName = pname
		item.ProductRating = rating
		item.ProductRetailers = Retailer.SmartPhoneRetailers(pname)
		item.ProductPrice = util.PriceStringTrimmer(price)
		item.ProductType = "smartphone"

		fmt.Println("img url is ->", imgUrl)
		// fmt.Println("hrefLink is ->", detailLink)
		fmt.Println("productname is ->", pname)
		fmt.Println("rating is ->", rating)
		fmt.Println("Retailer->", Retailer.SmartPhoneRetailers(pname))
		fmt.Println("global rating is ->", gRating)
		fmt.Println("price is ->", util.PriceStringTrimmer(price))

		productList = append(productList, item)
		fmt.Println()
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error while smartphone scraping->", err)
		anyErr = err
	})

	//we can iterate through all pages by looping.
	for i := 1; i <= 5; i++ {
		c.Visit(TARGET_LINK + strconv.Itoa(i))
	}
	c.Wait()
	fmt.Println("final post visit smartphone ")

	if anyErr != nil {
		return nil, anyErr
	}

	fmt.Println("the len is -", len(productList))

	return productList, nil
}
