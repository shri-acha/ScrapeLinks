package main

import (
	"fmt"
	// "io"
	"net/http"
	"golang.org/x/net/html"
)


type Link struct {
	linkURL string;
	isInternal bool;
} 

func main(){


	var link string
	fmt.Printf("Enter a (static) website to scrape for links:")
	

	fmt.Scanf("%s",&link)
	fmt.Printf("Starting the webscraper...\n")


	resp,err := http.Get(link)


	if err != nil{
		fmt.Printf("An Error Occured!: %v",err)
		return
	}

	defer resp.Body.Close()

	// body,err := io.ReadAll(resp.Body)
	

	z := html.NewTokenizer(resp.Body)

	for{

  if z.Next() == html.ErrorToken {
		return
	}
	tempByteString,_ := z.TagName()
		if string(tempByteString) == "a" {
		fmt.Printf("Result Body:%v\n",string(z.Raw()))
		}
	}
}
