package main

import (
	"strconv"
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"os"
	"bufio"
	"flag"
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

func (st ststats) calculateStockScore() int {
	score := 0

	// pegRatio	
	if (st.pegRatio > 0 && st.pegRatio < 1) {
		score++
	}
	if (st.pegRatio > 0 && st.pegRatio < 1.5) {
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

func readTickerFile(fn string) []string {
	var tickers []string
	ticre := regexp.MustCompile("\\(([A-Z]+)\\)")
    	//b, err := ioutil.ReadFile("shadow_stock_portfolio.csv")
    	b, err := ioutil.ReadFile(fn)
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
	
	//stock.score = calculateStockScore(stock)
	stock.score = stock.calculateStockScore()
	//stock.calculateStockScore() - how do I do this?
	
	return stock
}

func createStockCSV(l []ststats, fn string) {
	fo, err := os.Create(fn)
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

func worker(id int, tickers <-chan string, results chan<- ststats) {
    for t := range tickers {
        fmt.Println("worker", id, "processing ticker", t)
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
		results <- createStock(t,body)
    }
}

// interfaces
func main() {
	
	fName := flag.String("in", "", "grab the tickers from a file")
	oFile := flag.String("out", "results.csv", "name of the results file to be created")
	args := os.Args[1:]
	flag.Parse()
	var stocks []ststats
	var stock1 ststats
	fmt.Println(args)
	var newTickers []string

    tickers := make(chan string, 100)
    results := make(chan ststats, 100)

	for w := 1; w <= 3; w++ {
        go worker(w, tickers, results)
    }

	
	if (len(*fName) > 0) {	
		newTickers = readTickerFile(*fName)
		fmt.Println(newTickers)
	} else {
		newTickers = args
	}

	tcount := 0
	for _,t := range newTickers {
		
		tickers <- t
		tcount++
	}

	for a := 1; a <= tcount; a++ {
		stock1 = <-results
		stocks = append(stocks,stock1)
    }

	close(tickers)
	createStockCSV(stocks,*oFile)
	
}
