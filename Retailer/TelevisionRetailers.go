package Retailer

import (
	"math/rand"
	"time"
)

var televisionRetailers = []string{
	"https://vutvs.com/shop",
	"https://www.sony.co.in/bravia/gallery",
	"https://toshibatv-in.com/tvs-demo/",
	"https://www.samsung.com/in/tvs/all-tvs/",
	"https://acertvindia.com/product/",
	"https://www.mi.com/in/product-list/tv/",
	"https://www.lg.com/us/tvs",
	"https://www.panasonic.com/in/consumer/home-entertainment/televisions.html",
	"https://www.iffalcon.com/in/en/tvs",
	"https://www.vizio.com/en/shop/tv",
	"https://www.hisense-india.com/c/uhd-tv",
	"https://www.philips.co.in/c-m-so/tv",
	"https://www.haier.com/in/tvs/",
}

func TelevisionRetailers() string {
	rand.Seed(time.Now().UnixNano())
	var i int = rand.Intn(len(televisionRetailers))
	return televisionRetailers[i]
}
