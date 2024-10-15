package main

import (
	"fmt"
	"net/http"
	"strings"
	"golang.org/x/net/html"
)

type Link struct {
	linkURL    string
	isInternal bool
}

func main() {
	var link string
	fmt.Printf("Enter a (static)website to scrape for links: ")

	fmt.Scanf("%s", &link)
	fmt.Printf("Starting the webscraper...\n")

	resp, err := http.Get(link)
	if err != nil {
		fmt.Printf("An Error Occured!: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var tags []string
	tags = append(tags, "a")
	scrapedLinks := printScraped(tags, resp,link)	
	for i := range scrapedLinks {
		fmt.Printf("Link:%v\tIs-Internal:%v\n",scrapedLinks[i].linkURL,scrapedLinks[i].isInternal)
	}
}

func printScraped(links []string, resp *http.Response,link string) []Link{
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
					if len(tempValue) > 0 && (tempValue[0] == '/' || strings.Contains(link, tempValue)) {
						tempLinkArrBuff = append(tempLinkArrBuff, Link{tempValue, true})
					}else {
						tempLinkArrBuff = append(tempLinkArrBuff, Link{tempValue, false})
						}
					}
					if !moreAttr {
						break
					}	
				}
			}
		}
	}
	return tempLinkArrBuff
}
