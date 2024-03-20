package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

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

func main() {

	const LINK string = "https://www.amazon.in/s?i=electronics&bbn=976419031&rh=n%3A12045104031&ref=mega_elec_s23_3_1_1_5"

	const LINK1 string = "https://www.amazon.in/s?i=electronics&bbn=976419031&rh=n%3A12045104031&ref=sr_pg_"

	var imgList []string

	//wait group.
	var wg sync.WaitGroup //-> these helps to lock the Critical section during concurrent operation.

	c := colly.NewCollector(
	// colly.AllowedDomains("amazon.in"),
	)

	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("User-Agent", randomUserAgent()) //->this is done to prevent the browser letting scraper to block.
	})

	c.OnHTML("div.s-result-item", func(h *colly.HTMLElement) {
		imgUrl := h.ChildAttr("img.s-image", "src")

		if len(imgUrl) != 0 {
			imgList = append(imgList, imgUrl)
		}

	})

	wg.Add(1)

	//use of go-routines that are useful to carry out async tasks.
	go func() {
		for i := 1; i < 10; i++ {

			err := c.Visit(LINK1 + strconv.Itoa(i))

			if err != nil {
				log.Fatal("The error is ", err)
			}

		}
		defer wg.Done()
	}()

	wg.Wait()

	for i, val := range imgList {
		fmt.Println("The img link is ", val, "<-", i)
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
