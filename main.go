package main

import (
	"fmt"
	"log"

	// "strconv"
	"strings"
	"sync"
	"time"

	model "github.com/Ivan2001otp/Visionary-AI/Model"
	"github.com/Ivan2001otp/Visionary-AI/Retailer"
	scrap "github.com/Ivan2001otp/Visionary-AI/Service/Scrap"

	// config "github.com/Ivan2001otp/Visionary-AI/Service/Database"
	util "github.com/Ivan2001otp/Visionary-AI/Util"
	"github.com/gocolly/colly"
)

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

func display(str string) {
	fmt.Println(str)
}

func isEmpty(str string) bool {
	return len(str) == 0
}

func isProductEmpty(product model.Product) bool {
	if isEmpty(product.ProductName) && isEmpty(product.ProductDetailUrl) && isEmpty(product.ProductRating) && isEmpty(product.ProductImg) {
		return true
	}
	return false
}

func main() {
	// var i int
	const LINK string = "https://www.amazon.in/s?i=electronics&bbn=976419031&rh=n%3A12045104031&ref=mega_elec_s23_3_1_1_5"

	const LINK1 string = "https://www.amazon.in/s?i=electronics&bbn=976419031&rh=n%3A12045104031&ref=sr_pg_"

	// var imgList []string
	// var href []string
	// var pname []string
	var productList []model.Product
	//wait group.
	var wg sync.WaitGroup //-> these helps to lock the Critical section during concurrent operation.

	c := colly.NewCollector(
		// colly.AllowedDomains("amazon.in"),
		colly.UserAgent(util.RandomUserAgent()),
	)

	// c.UserAgent = randomUserAgent()

	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("User-Agent", randomUserAgent()) //->this is done to prevent the browser letting scraper to block.

	})

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

		if !isEmpty(price) {
			productClone.ProductPrice = priceStringTrimmer(price)
		}

		if len(imgUrl) != 0 {
			// imgList = append(imgList, imgUrl)
			productClone.ProductImg = imgUrl
		}

		if len(hrefLink) != 0 {

			if !strings.HasPrefix(hrefLink, "https://www.amazon.in") {
				hrefLink = "https://www.amazon.in" + hrefLink
			}
			productClone.ProductDetailUrl = hrefLink
		}

		if !isEmpty(pname) {
			productClone.ProductName = pname
		}

		if !isEmpty(rating) {
			productClone.ProductRating = rating
		}

		if !isEmpty(globalRating_) {

			productClone.GlobalRating = globalRating_
		}

		productClone.ProductRetailers = Retailer.TelevisionRetailers()

		if !isProductEmpty(productClone) {
			productList = append(productList, productClone)
		}

	})

	wg.Add(1)
	startTime := time.Now()
	/*
		//use of go-routines that are useful to carry out async tasks.
		// This goroutine service fetches the tele info
		go func() {

			for i := 1; i < 2; i++ {

				err := c.Visit(LINK1 + strconv.Itoa(i))

				if err != nil {
					log.Fatal("The error is ", err)
				}

			}

			defer wg.Done()
		}()
	*/

	//this go routine fetches smartphone info.
	go func() {

		defer wg.Done()
		res, err := scrap.WatchScraper()

		if err != nil {
			log.Fatal(err)
		}

		for i, item := range res {
			result, err := item.SaveToMongo()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Inserted ", i+1, "->", result)
		}
		fmt.Println("the len is ", len(res))
	}()

	// scraps the info about laptops
	/*
		go func() {
			defer wg.Done()
			productList, err := scrap.LapTopScraper()

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("The response has ->", len(productList))
		}()
	*/
	/*
		//database service...

		go func() {
			defer wg.Done()

			mongoClient, err := config.MongoDbInstanceProvider()

			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Successfully connected!", mongoClient)

			}
			if mongoClient != nil {
				var p model.Product
				p.ProductName = "lg brand watch"
				p.CategoryType = "Electronics"
				p.GlobalRating = "3"
				p.ProductDetailUrl = "href.img"
				p.ProductImg = "abc.jpeg"
				p.ProductPrice = "1230"
				p.ProductRating = "4"
				p.ProductRetailers = "hi bro Inc"

				res, err := p.SaveToMongo()
				if err != nil {
					log.Fatal("error on saving data from driver : ", err)
				}

				fmt.Println("The data is saved successfully ->", res)
			}
		}()
	*/
	wg.Wait()

	endTime := time.Now()
	fmt.Println("The time taken is ", endTime.Sub(startTime))

	for i, val := range productList {
		fmt.Println(i+1, "rating->", val.ProductRating)
		fmt.Println(i+1, "name->", val.ProductName)
		// fmt.Println(i+1, "href->", val.productDetailUrl)
		fmt.Println(i+1, "price->", val.ProductPrice)
		fmt.Println(i+1, "imgUrl->", val.ProductImg)
		fmt.Println(i+1, "globalRating->", val.GlobalRating)
		fmt.Println(i+1, "retailer->", val.ProductRetailers)
		fmt.Println()
	}

}
