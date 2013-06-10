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

// populate fields in structure for one ticker
// pass tickers in via command line
// loop through, populating the full structure
// score each ticker
func main() {
        argsWithoutProg := os.Args[1:]
        fmt.Println(argsWithoutProg)
        // Profit Margin (ttm):</td><td class="yfnc_tabledata1">21.58%</td>
        //r := regexp.MustCompile("Profit Margin (ttm):</td><td [^>].>([^<]+)</td>")
        pmre := regexp.MustCompile("Profit Margin.*?</td><td.*?>(.*?)</td>")
        prre := regexp.MustCompile("PEG Ratio.*?</td><td.*?>(.*?)</td>")



        //res, err := http.Get("http://finance.yahoo.com/d/quotes.csv?s=bwld+ctsh+cybx+dar+jazz+mwiv+pcp+pets+rex+voxx&f=snl1")
        res, err := http.Get("http://finance.yahoo.com/q/ks?s=msft")
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
        stock1 := ststats{"msft",0,0,0,0,0,0}

        strs := pmre.FindSubmatch(body)
	stock1.pMargin = convertToFloat(string(strs[1]))
	
        strs = prre.FindSubmatch(body)
	stock1.pegRatio = convertToFloat(string(strs[1]))
	
	
	fmt.Println(stock1)
	/*for _, line := range strings.Split(string(body), "\n") {
		if (len(line) > 0) {
			// TBD - parse CSV line into fields, store in DB
			fmt.Println("line splitter =",line)
	  	}
	} */ 

}
