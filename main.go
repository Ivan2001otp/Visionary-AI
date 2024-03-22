package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Ivan2001otp/Visionary-AI/Retailer"
	"github.com/gocolly/colly"
)

type Product struct {
	productName      string
	productImg       string
	productDetailUrl string
	productRating    string
	globalRating     string
	productPrice     string
	productRetailers string
}

var UserAgentStrings = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.67 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:98.0) Gecko/20100101 Firefox/98.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/103.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Opera/93.0.4585.67 Safari/537.36",
}

func randomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	var i int = rand.Intn(len(UserAgentStrings))
	return UserAgentStrings[i]
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

func display(str string) {
	fmt.Println(str)
}

func isEmpty(str string) bool {
	return len(str) == 0
}

func isProductEmpty(product Product) bool {
	if isEmpty(product.productName) && isEmpty(product.productDetailUrl) && isEmpty(product.productRating) && isEmpty(product.productImg) {
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
	var productList []Product

	//wait group.
	var wg sync.WaitGroup //-> these helps to lock the Critical section during concurrent operation.

	c := colly.NewCollector(
	// colly.AllowedDomains("amazon.in"),
	)

	// c.UserAgent = randomUserAgent()

	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("User-Agent", randomUserAgent()) //->this is done to prevent the browser letting scraper to block.

	})

	c.OnHTML("div.s-result-item", func(h *colly.HTMLElement) {

		// p := h.Response
		// i += 1
		// fmt.Println(i, "->", p.StatusCode)
		var productClone Product

		imgUrl := h.ChildAttr("img.s-image", "src")                        //img
		hrefLink := h.ChildAttr("a.a-link-normal", "href")                 //href
		pname := h.ChildText("a.a-link-normal > span.a-text-normal")       // pname
		rating := h.ChildText("i.a-icon-star-small > span.a-icon-alt")     //rating out 5.
		price := h.ChildText("span.a-price > span.a-offscreen")            //price.
		globalRating_ := h.ChildText("a.a-link-normal > span.a-size-base") //global rating

		//a.a-link-normal > span.a-price-whole
		if !isEmpty(price) {
			productClone.productPrice = priceStringTrimmer(price)
		}

		// p := h.Attr("span.a-size-base-plus a-color-base a-text-normal")
		if len(imgUrl) != 0 {
			// imgList = append(imgList, imgUrl)
			productClone.productImg = imgUrl
		}

		if len(hrefLink) != 0 {

			if !strings.HasPrefix(hrefLink, "https://www.amazon.in") {
				hrefLink = "https://www.amazon.in" + hrefLink
			}
			// href = append(href, hrefLink)
			productClone.productDetailUrl = hrefLink
		}

		if !isEmpty(pname) {
			// pname = append(pname, p)
			productClone.productName = pname
		}

		if !isEmpty(rating) {
			productClone.productRating = rating
		}
		if !isEmpty(globalRating_) {

			// gRating, err := strconv.Atoi(globalRating_)

			// if err != nil {
			// 	log.Fatal("Something went wrong while parsing GlobalRating from String to int32:", err)
			// }
			productClone.globalRating = globalRating_
		}
		productClone.productRetailers = Retailer.TelevisionRetailers()
		if !isProductEmpty(productClone) {
			productList = append(productList, productClone)
		}

	})

	wg.Add(1)
	startTime := time.Now()
	//use of go-routines that are useful to carry out async tasks.
	go func() {

		for i := 1; i < 2; i++ {

			err := c.Visit(LINK1 + strconv.Itoa(i))

			if err != nil {
				log.Fatal("The error is ", err)
			}

		}

		defer wg.Done()
	}()

	wg.Wait()

	// for i, val := range imgList {
	// 	fmt.Println("The img link is ", val, "<-", i)

	// }
	endTime := time.Now()
	fmt.Println("The time taken is ", endTime.Sub(startTime))

	for i, val := range productList {
		fmt.Println(i+1, "rating->", val.productRating)
		fmt.Println(i+1, "name->", val.productName)
		// fmt.Println(i+1, "href->", val.productDetailUrl)
		fmt.Println(i+1, "price->", val.productPrice)
		fmt.Println(i+1, "imgUrl->", val.productImg)
		fmt.Println(i+1, "globalRating->", val.globalRating)
		fmt.Println(i+1, "retailer->", val.productRetailers)
		fmt.Println()
	}

}

/*

Web scraping using ChromeDp - headless browser.
func main() {

	const LINK = "https://www.amazon.in/s?i=electronics&bbn=976419031&rh=n%3A12045104031&ref=mega_elec_s23_3_1_1_5"
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)

	//scrolling code in JS
	scrollingScript := `
	// scroll down the page 8 times
	const scrolls = 30
	let scrollCount = 0

	// scroll down and then wait for 0.5s
	const scrollInterval = setInterval(() => {
	  window.scrollTo(0, document.body.scrollHeight)
	  scrollCount++

	  if (scrollCount === numScrolls) {
	   clearInterval(scrollInterval)
	  }
	}, 500)
 `

	defer cancel()

	var imgNodes []*cdp.Node
	var url string
	var isTrue bool
	var pName string
	err := chromedp.Run(
		ctx,
		chromedp.Navigate(LINK),
		// chromedp.Sleep(3*time.Second),
		//.puis-card-container

		chromedp.Evaluate(scrollingScript, nil),
		chromedp.WaitVisible("div.s-result-item"),

		//.s-main-slot s-result-list s-search-results
		chromedp.Nodes("div.s-result-item", &imgNodes, chromedp.ByQueryAll),
	)

	if err != nil {
		log.Fatal("Error while performing the automation logic:", err)
	}

	fmt.Println("The len is ", len(imgNodes))

	for _, node := range imgNodes {
		err = chromedp.Run(ctx,
			chromedp.AttributeValue("img.s-image", "src", &url, &isTrue, chromedp.FromNode(node)),
			chromedp.AttributeValue("img.s-image", "alt", &pName, &isTrue, chromedp.FromNode(node)),
		)

		if err != nil {
			log.Fatal("The error is ", err)
		}

		fmt.Println("THe link is ", url)
		fmt.Println("The p-name is ", pName)
		fmt.Println()

	}
}
*/
