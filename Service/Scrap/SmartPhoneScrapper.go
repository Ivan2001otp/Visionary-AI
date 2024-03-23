package scrap

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Ivan2001otp/Visionary-AI/Retailer"
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

func priceStringTrimmer(str string) string {
	const rupaySymbol = '₹'
	var temp string = ""
	var rupayFreq int = 0
	for _, char := range str {
		if char == rupaySymbol && rupayFreq != 1 {
			rupayFreq++
			temp = temp + string(char) //append first ruppee symbol
		} else if char != rupaySymbol && rupayFreq == 1 {
			temp = temp + string(char)
		} else if char == rupaySymbol && rupayFreq == 1 {
			break
		}
	}

	return temp
}

func SmartPhoneScrapper() ([]SmartPhone, error) {
	fmt.Println("visiting smartphone ")
	var productList []SmartPhone
	const TARGET_LINK = "https://www.amazon.in/s?k=smartphones&i=electronics&rh=n%3A1389401031&page=2&ref=sr_pg_"
	var anyErr error

	c := colly.NewCollector()

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
		item.ProductPrice = priceStringTrimmer(price)
		item.ProductType = "smartphone"

		fmt.Println("img url is ->", imgUrl)
		// fmt.Println("hrefLink is ->", detailLink)
		fmt.Println("productname is ->", pname)
		fmt.Println("rating is ->", rating)
		fmt.Println("Retailer->", Retailer.SmartPhoneRetailers(pname))
		fmt.Println("global rating is ->", gRating)
		fmt.Println("price is ->", priceStringTrimmer(price))

		productList = append(productList, item)
		fmt.Println()
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error while smartphone scraping->", err)
		anyErr = err
	})

	//we can iterate through all pages by looping.
	for i := 1; i < 2; i++ {
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
