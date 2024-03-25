package util

import (
	"fmt"
	"math/rand"
	"time"

	model "github.com/Ivan2001otp/Visionary-AI/Model"
)

var userAgentStrings = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.67 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:98.0) Gecko/20100101 Firefox/98.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/103.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Opera/93.0.4585.67 Safari/537.36",
}

func RandomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	var i int = rand.Intn(len(userAgentStrings))
	return userAgentStrings[i]
}

func IsEmpty(str string) bool {
	return len(str) == 0
}

func Print(str string) {
	fmt.Println(str)
}

func IsProductInfoEmpty(product model.Product) bool {
	if IsEmpty(product.ProductName) && IsEmpty(product.ProductDetailUrl) && IsEmpty(product.ProductRating) && IsEmpty(product.ProductImg) && IsEmpty(product.GlobalRating) && IsEmpty(product.ProductPrice) {
		return true
	}
	return false
}

func PriceStringTrimmer(str string) string {
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
