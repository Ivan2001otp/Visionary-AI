package util

import (
	"fmt"

	model "github.com/Ivan2001otp/Visionary-AI/Model"
)

func IsEmpty(str string) bool {
	return len(str) == 0
}

func Print(str string) {
	fmt.Println(str)
}

func IsProductEmpty(product model.Product) bool {
	if IsEmpty(product.ProductName) && IsEmpty(product.ProductDetailUrl) && IsEmpty(product.ProductRating) && IsEmpty(product.ProductImg) {
		return true
	}
	return false
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
