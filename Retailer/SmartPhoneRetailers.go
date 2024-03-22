package Retailer

import "strings"

func SmartPhoneRetailers(str string) string {
	retailerSites := map[string]string{
		"apple":    "https://www.apple.com/in/iphone/",
		"samsung":  "https://www.samsung.com/in/smartphones/all-smartphones/",
		"oneplus":  "https://www.oneplus.in/store/phone",
		"redmi":    "https://www.mi.com/global/product-list/redmi/",
		"xiaomi":   "https://www.mi.com/in/product-list/xiaomi/",
		"realme":   "https://www.realme.com/in/realme-phones",
		"motorola": "https://www.motorola.in/smartphones",
		"poco":     "https://www.poco.in/",
		"nokia":    "https://www.hmd.com/en_in/smartphones",
		"nothing":  "https://in.nothing.tech/",
		"mi":       "https://www.mi.com/in/",
		"sony":     "https://electronics.sony.com/c/mobile",
	}

	var temp string = ""

	for _, val := range str {
		if val != ' ' {
			temp = temp + string(val)
		} else {
			temp = strings.ToLower(temp)
			break
		}
	}
	if retailerSites[temp] == "" {
		return "https://www.alibaba.com/"
	}
	return retailerSites[temp]
}
