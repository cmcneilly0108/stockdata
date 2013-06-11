package main

import (
	"strconv"
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"os"
	)

type ststats struct {
	ticker string
	score int
	pegRatio float64
	pMargin float64
	revenueGrowth float64
	yrChange float64
	spyrChange float64
	netChange float64
}

func calculateStockScore(st ststats) int {
	score := 0
	if (st.pegRatio < 1) {
		score++
	}
	if (st.pegRatio < 1.5) {
		score++
	}
	if (st.pegRatio > 2.5) {
		score--
	}
	if (st.pMargin > .1) {
		score++
	}
	if (st.pMargin > .3) {
		score++
	}
	if (st.pMargin < 0) {
		score--
	}
	if (st.revenueGrowth > .1) {
		score++
	}
	if (st.revenueGrowth > .3) {
		score++
	}
	if (st.revenueGrowth < 0) {
		score--
	}
	if (st.netChange > 0) {
		score++
	}
	if (st.netChange > .1) {
		score++
	}
	if (st.netChange < 0) {
		score--
	}
	
	return score
}

func convertToFloat(numString string) float64 {
	var result float64
	if numString[len(numString)-1] == '%' {
		numString = numString[0:len(numString)-1]
		result, _ = strconv.ParseFloat(numString,64)
		result = result/100
	} else {
		result, _ = strconv.ParseFloat(numString,64)
	}
	return result
}

// pass tickers in via command line
// loop through, populating the full structure for all tickers
func main() {
	argsWithoutProg := os.Args[1:]
	fmt.Println(argsWithoutProg)
	ticker := argsWithoutProg[0]
	url := "http://finance.yahoo.com/q/ks?s=" + ticker

	pmre := regexp.MustCompile("Profit Margin.*?</td><td.*?>(.*?)</td>")
	prre := regexp.MustCompile("PEG Ratio.*?</td><td.*?>(.*?)</td>")
	rgre := regexp.MustCompile("Qtrly Revenue.*?</td><td.*?>(.*?)</td>")
	ycre := regexp.MustCompile("52-Week Change.*?</td><td.*?>(.*?)</td>")
	scre := regexp.MustCompile("P500 52-Week Change.*?</td><td.*?>(.*?)</td>")
	
	
	
	//res, err := http.Get("http://finance.yahoo.com/d/quotes.csv?s=bwld+ctsh+cybx+dar+jazz+mwiv+pcp+pets+rex+voxx&f=snl1")
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get", err)
		return
	}
	defer res.Body.Close()
	//Grab the results from the call into a byte array?
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll", err)
		return
	}
	stock1 := ststats{ticker,0,0,0,0,0,0,0}
	
	strs := pmre.FindSubmatch(body)
	stock1.pMargin = convertToFloat(string(strs[1]))
	
	strs = prre.FindSubmatch(body)
	stock1.pegRatio = convertToFloat(string(strs[1]))
	
	strs = rgre.FindSubmatch(body)
	stock1.revenueGrowth = convertToFloat(string(strs[1]))
	
	strs = ycre.FindSubmatch(body)
	stock1.yrChange = convertToFloat(string(strs[1]))
	
	strs = scre.FindSubmatch(body)
	stock1.spyrChange = convertToFloat(string(strs[1]))
	
	stock1.netChange = stock1.yrChange - stock1.spyrChange
	
	stock1.score = calculateStockScore(stock1)
	
	
	fmt.Println(stock1)
	
}
