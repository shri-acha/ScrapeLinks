package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"golang.org/x/net/html"
)

type Link struct {
	linkURL    string
	content string
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
				fmt.Printf("Meta: %v\tContent: %v\n", link.linkURL, link.content)
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
	tags = append(tags, "meta")
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
          var tempValue, content string
          if string(key) == "property" {
            tempValue = string(value)
            if moreAttr {
              key,value,_ := z.TagAttr();
              if string(key) == "content"{
                content = string(value)
              }           }
           tempLinkArrBuff = append(tempLinkArrBuff, Link{tempValue, content})
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
