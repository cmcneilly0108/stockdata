package main

import (
        "fmt"
        "net/http"
        "io/ioutil"
        "regexp"
//        "strings"
)

// http://finance.yahoo.com/d/quotes.csv?s=bwld+ctsh&f=snl1

func main() {
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
        //Convert the byte array into each line, and convert to a string
        //fmt.Println(string(body))
        // Profit Margin (ttm):</td><td class="yfnc_tabledata1">21.58%</td>
        //r := regexp.MustCompile("Profit Margin (ttm):</td><td [^>].>([^<]+)</td>")
        r := regexp.MustCompile("Profit Margin[^<]+</td><td [^>]+>([^<]+)")
        strs := r.FindSubmatch(body)
        fmt.Println(string(strs[1]))
	/*for _, line := range strings.Split(string(body), "\n") {
		if (len(line) > 0) {
			// TBD - parse CSV line into fields, store in DB
			fmt.Println("line splitter =",line)
	  	}
	} */ 

}
