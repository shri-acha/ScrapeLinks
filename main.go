package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"golang.org/x/net/html"
)

type Link struct {
	linkURL    string
	isInternal bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "At least 1 argument is required!\nSyntax: ./main -u [URL]")
		return
	}
	if (os.Args[1] == "-u") && (len(os.Args[2]) > 0) {
		links := os.Args[2:]
		var wg sync.WaitGroup
		results := make(chan []Link, len(links))
		for _, link := range links {
			wg.Add(1)
			go func(link string) {
				defer wg.Done()
				printScraped(link, results)
			}(link)
		}
		go func() {
			wg.Wait()
			close(results)
		}()
		for scrapedLinks := range results {
			for _, link := range scrapedLinks {
				fmt.Printf("Link: %v\tIs-Internal: %v\n", link.linkURL, link.isInternal)
			}
		}
	} else {
		fmt.Printf("Syntax: ./main -u [LINK]")
	}
}

func printScraped(link string, results chan<- []Link) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Printf("An error occurred while fetching %s: %v\n", link, err)
		results <- []Link{}
		return
	}
	defer resp.Body.Close()
	var tags []string
	tags = append(tags, "a")
	fmt.Println("Scraped links for the website:\t" + link)
	scrapeLinks(tags, resp, link, results)
}

func scrapeLinks(links []string, resp *http.Response, baseLink string, results chan<- []Link) {
	z := html.NewTokenizer(resp.Body)
	var tempLinkArrBuff []Link
	for {
		tokenType := z.Next()
		if tokenType == html.ErrorToken {
			break
		}
		tempByteString, _ := z.TagName()
		for _, linkTag := range links {
			if string(tempByteString) == linkTag {
				for {
					key, value, moreAttr := z.TagAttr()
					if string(key) == "href" {
						tempValue := string(value)
						isInternal := len(tempValue) > 0 && (tempValue[0] == '/' || strings.Contains(tempValue, baseLink))
						tempLinkArrBuff = append(tempLinkArrBuff, Link{tempValue, isInternal})
					}
					if !moreAttr {
						break
					}
				}
			}
		}
	}
	results <- tempLinkArrBuff
}
