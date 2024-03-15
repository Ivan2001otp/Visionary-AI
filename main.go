package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	// "github.com/PuerkitoBio/goquery"
)

var userAgents = []string{
	// "Mozilla/5.0 (Windows NT 10.0; Win64: x64) AppleWebKit/537.36 (KHTML, like Gecko/61.0.31)",
	// "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/100.0 Mobile Safari/537.36",
	// "Mozilla/5.0 (Windows NT 11.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/105.0.0.0 Safari/537.36",
	// "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0",

	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
}

const LIST_URL string = "div.organic-list organic-list_G app-organic-search__list organic-list-gallery-wrapper"

type EcommerceProduct struct {
	productName     string
	productImageUrl string
	productPrice    string
	productSupplier string
	monthlySales    string
	ratings         string
}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())

	randNumber := rand.Int() % len(userAgents)

	fmt.Println("The user-agent is ", userAgents[randNumber])
	return userAgents[randNumber]
}

/*
func getResponseBody(targetUrl string) (*http.Response, error) {

		client := &http.Client{}

		req, err := http.NewRequest("GET", targetUrl, nil)

		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", randomUserAgent())
		res, err := client.Do(req)

		if err != nil {
			return nil, err
		}

		//_, _ := io.ReadAll(res.Body)

		//fmt.Println(string(r))
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			fmt.Println("Bad request")
		}

		return res, nil


}
*/

/*
func scrapData(baseUrl string, response *http.Response) []string {
	//fetches images.

	if response != nil {
		//here u try to fetch the html data.
		doc, _ := goquery.NewDocumentFromReader(response.Body)
		foundImgUrls := []string{}

		if doc != nil {

			//search-card-e-slider__wrapper
			//a
			//img->src
			// doc.Filter(LIST_URL).Each(func(i1 int, s1 *goquery.Selection) {
			// 	s1.Find("img").Each(func(i2 int, s2 *goquery.Selection) {
			// 		x, _ := s2.Attr("src")
			// 		if !strings.HasPrefix(x, "https:") {
			// 			x = "https:" + x
			// 		}
			// 		fmt.Println("Img link -> ", x)
			// 		foundImgUrls = append(foundImgUrls, x)
			// 	})
			// })

			doc.Find("img.search-card-e-slider__img").Each(func(i int, s *goquery.Selection) {
				x, _ := s.Attr("src")
				if !strings.HasPrefix(x, "https:") {
					x = "https:" + x
				}
				fmt.Println("Img link -> ", x)

				foundImgUrls = append(foundImgUrls, x)
			})

			return foundImgUrls
		}
	}

	return []string{"Nothing"}

}
*/

func Crawl(baseUrl string, ctx context.Context) []string {

	fmt.Println(baseUrl)

	var imageNodes []*cdp.Node

	var imageList []string
	var temp string
	var ok bool

	err := chromedp.Run(ctx,
		//visit the target url
		chromedp.Navigate(baseUrl),

		//wait for the page to load.
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		// chromedp.AttributeValue(`img`, "search-card-e-slider__img", &temp, &isPresent),
		// imageNodes = append(imageNodes,temp);
		chromedp.Nodes(".search-card-e-slider__img", &imageNodes, chromedp.ByQueryAll),
	)

	if err != nil {
		fmt.Println("Error captured - first!")
		log.Fatal(err)
	}

	//fetch from nodes.
	for _, node := range imageNodes {
		//extract data from imageNode HTML
		err = chromedp.Run(ctx, chromedp.AttributeValue(node, "src", &temp, &ok, nil))

		if err != nil {
			fmt.Println("Error captured while iterating..")
			log.Fatal(err)
		}

		imageList = append(imageList, temp)
	}
	return imageList

	/*
		requestBody, err := getResponseBody(baseUrl)

		// fmt.Println(string(r))
		if err != nil {
			panic(err)
		}

		//links := scrapData(baseUrl, requestBody)

		return links*/

}

func main() {
	// baseDomain := "https://www.alibaba.com/trade/search?spm=a2700.galleryofferlist.pageModule_fy23_pc_search_bar.keydown__Enter&tab=all&searchText=snickers+shoes+for+men"
	baseDomain := "https://www.alibaba.com/trade/search?IndexArea=product_en&CatId=&fsb=y&viewtype=&tab=all&SearchScene=&SearchText=clothes+for+man+and+woman"

	options := []chromedp.ExecAllocatorOption{
		// chromedp.DefaultExecAllocatorOptions[0],
		chromedp.ProxyServer("50.206.111.88:80"),
		chromedp.UserAgent(randomUserAgent()),
	}
	context, cancel := chromedp.NewExecAllocator(
		context.Background(),
		options...,
	)
	defer cancel()

	context1, cancel1 := chromedp.NewContext(context)
	defer cancel1()

	var wg sync.WaitGroup
	wg.Add(1)

	go func(baseUrl string, wg *sync.WaitGroup) {

		var itemCollectionLinks []string

		for len(itemCollectionLinks) <= 1 {

			itemCollectionLinks = Crawl(baseUrl, context1)
			fmt.Println("The length is ", len(itemCollectionLinks))
			fmt.Println()

		}
		wg.Done()
	}(baseDomain, &wg)

	wg.Wait()

	fmt.Println("Finished - terminated")
}
