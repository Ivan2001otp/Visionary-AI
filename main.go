package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	model "github.com/Ivan2001otp/Visionary-AI/Model"
	"github.com/Ivan2001otp/Visionary-AI/Retailer"
	scrap "github.com/Ivan2001otp/Visionary-AI/Service/Scrap"
	"github.com/gocolly/colly"
)

var UserAgentStrings = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.67 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:98.0) Gecko/20100101 Firefox/98.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/103.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Opera/93.0.4585.67 Safari/537.36",
}

//	type Product struct {
//		productName      string
//		productImg       string
//		productDetailUrl string
//		productRating    string
//		globalRating     string
//		productPrice     string
//		productRetailers string
//	}
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

	var laptopList_ []scrap.LapTop

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
	//use of go-routines that are useful to carry out async tasks.
	/*This goroutine service fetches the tele info
	go func() {

		for i := 1; i < 2; i++ {

			err := c.Visit(LINK1 + strconv.Itoa(i))

			if err != nil {
				log.Fatal("The error is ", err)
			}

		}

		defer wg.Done()
	}()*/

	//this go routine fetches smartphone info.
	go func() {

		defer wg.Done()
		res, err := scrap.SmartPhoneScrapper()

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(len(res))
	}()

	/* scraps the info about laptops
	go func() {
		defer wg.Done()

		laptopList_, err := scrap.LapTopScraper()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(laptopList_)
	}()*/

	wg.Wait()

	endTime := time.Now()
	fmt.Println("The time taken is ", endTime.Sub(startTime))

	for i, val := range laptopList_ {
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
