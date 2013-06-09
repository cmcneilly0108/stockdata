package main

import (
        "fmt"
        "net/http"
        "io/ioutil"
)

// http://finance.yahoo.com/d/quotes.csv?s=bwld+ctsh&f=snl1

func main() {
        res, err := http.Get("http://finance.yahoo.com/d/quotes.csv?s=bwld+ctsh+cybx+dar+jazz+mwiv+pcp+pets+rex+voxx&f=snl1")
        if err != nil {
                fmt.Println("http.Get", err)
                return
        }
        defer res.Body.Close()
        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
                fmt.Println("ioutil.ReadAll", err)
                return
        }
        lenp := len(body)
        //if maxp := 60; lenp > maxp {
        //        lenp = maxp
        //}
        fmt.Println(len(body), string(body[:lenp]))
}
