package main

import (
	"strconv"
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"os"
	"bufio"
	)

var (
	pmre = regexp.MustCompile("Profit Margin.*?</td><td.*?>(.*?)</td>")
	prre = regexp.MustCompile("PEG Ratio.*?</td><td.*?>(.*?)</td>")
	rgre = regexp.MustCompile("Qtrly Revenue.*?</td><td.*?>(.*?)</td>")
	ycre = regexp.MustCompile("52-Week Change.*?</td><td.*?>(.*?)</td>")
	scre = regexp.MustCompile("P500 52-Week Change.*?</td><td.*?>(.*?)</td>")
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

	// pegRatio	
	if (st.pegRatio < 1) {
		score++
	}
	if (st.pegRatio < 1.5) {
		score++
	}
	if (st.pegRatio > 2.5) {
		score--
	}

	// Profit Margin
	if (st.pMargin > .1) {
		score++
	}
	if (st.pMargin > .3) {
		score++
	}
	if (st.pMargin < 0) {
		score--
	}

	// Revenue Growth
	if (st.revenueGrowth > .1) {
		score++
	}
	if (st.revenueGrowth > .3) {
		score++
	}
	if (st.revenueGrowth > .7) {
		score++
	}
	if (st.revenueGrowth < 0) {
		score--
	}
	
	// Net Change
	if (st.netChange > 0) {
		score++
	}
	if (st.netChange > .2) {
		score++
	}
	if (st.netChange < 0) {
		score--
	}
	if (st.netChange < -.1) {
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

func readTickerFile() []string {
	var tickers []string
	ticre := regexp.MustCompile("\\(([A-Z]+)\\)")
    	b, err := ioutil.ReadFile("shadow_stock_portfolio.csv")
    	if err != nil { panic(err) }
    	
	fstr := string(b)
	t2 := ticre.FindAllStringSubmatch(fstr, -1)
	fmt.Println(t2)
	for _,s := range t2 {
		tickers = append(tickers,s[1])
	}
	return tickers
}

func createStock(ticker string, body []byte) ststats {
	stock := ststats{ticker,0,0,0,0,0,0,0}
	
	strs := pmre.FindSubmatch(body)
	stock.pMargin = convertToFloat(string(strs[1]))
	
	strs = prre.FindSubmatch(body)
	stock.pegRatio = convertToFloat(string(strs[1]))
	
	strs = rgre.FindSubmatch(body)
	stock.revenueGrowth = convertToFloat(string(strs[1]))
	
	strs = ycre.FindSubmatch(body)
	stock.yrChange = convertToFloat(string(strs[1]))
	
	strs = scre.FindSubmatch(body)
	stock.spyrChange = convertToFloat(string(strs[1]))
	
	stock.netChange = stock.yrChange - stock.spyrChange
	
	stock.score = calculateStockScore(stock)
	
	return stock
}

func createStockCSV(l []ststats) {
	fo, err := os.Create("stocks.csv")
	if err != nil { panic(err) }

	// make a write buffer
    	w := bufio.NewWriter(fo)
    	header := "Ticker,Score,PEG Ratio,Profit Margin,YOY Growth,Growth Diff\n"
	if _, err := w.WriteString(header); err != nil {
	    panic(err)
	}

    	for _,st := range l {
    		line := st.ticker + "," + strconv.Itoa(st.score) + ","
    		s := strconv.FormatFloat(st.pegRatio,'f',2,64)
    		line +=  s + ","
    		s = strconv.FormatFloat(st.pMargin,'f',3,64)
    		line +=  s + ","
    		s = strconv.FormatFloat(st.revenueGrowth,'f',3,64)
    		line +=  s + ","
    		s = strconv.FormatFloat(st.netChange,'f',3,64)
    		line +=  s
    		line +=  "\n"
		// write a chunk
		if _, err := w.WriteString(line); err != nil {
		    panic(err)
		}
	}

	if err = w.Flush(); err != nil { panic(err) }
	fo.Close()
}

// pass command line flag to process args or file - file name

func main() {
	args := os.Args[1:]
	var stocks []ststats
	fmt.Println(args)
	
	newTickers := readTickerFile()
	fmt.Println(newTickers)

	for _,t := range newTickers {
		url := "http://finance.yahoo.com/q/ks?s=" + string(t)
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("http.Get", err)
			return
		}

		//Grab the results from the call into a byte array?
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("ioutil.ReadAll", err)
			return
		}
		res.Body.Close()
		stock1 := createStock(t,body)
		stocks = append(stocks,stock1)
		
		fmt.Println(stock1)
	}
	createStockCSV(stocks)
	
}
