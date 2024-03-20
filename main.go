package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {

	const LINK = "https://www.amazon.in/s?i=electronics&bbn=976419031&rh=n%3A12045104031&ref=mega_elec_s23_3_1_1_5"
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)

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
