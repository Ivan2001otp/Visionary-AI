package Retailer

import (
	"strings"
)

func WatchRetailer(key string) string {
	const defaultRetailer string = "https://www.myntra.com/smart-watches"

	var retailerSites = map[string]string{
		"fastrack":   "https://www.fastrack.in/shop/watches",
		"noise":      "https://www.gonoise.com/collections/smart-watches",
		"apple":      "https://www.apple.com/in/shop/buy-watch/apple-watch",
		"boat":       "https://www.boat-lifestyle.com/collections/smart-watches",
		"fire-boltt": "https://www.fireboltt.com/",
		"google":     "https://store.google.com/product/pixel_watch_2?hl=en-GB&pli=1",
	}

	var temp string = ""
	for i, _ := range key {
		if key[i] != ' ' {
			temp = temp + string(key[i])
		} else {
			temp = strings.ToLower(temp)
		}
	}

	if retailerSites[temp] == "" {
		return defaultRetailer
	}
	return retailerSites[temp]

}
