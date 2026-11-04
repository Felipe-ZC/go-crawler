package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const ProtocolStr = "https://"

func main() {
	client := &http.Client{Timeout: 30 * time.Second}
	
	fmt.Println("Fetching...")
	baseUrl := "https://github.com"

	resp, err := client.Get(baseUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	
	fmt.Println("Status code: ", resp.StatusCode)
	if resp.StatusCode >= 400 {
		fmt.Println("Non-200 status code returned, exiting...")
		return
	}

	fmt.Println("Parsing body...")
	doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}
	// TODO: Find all URLs here!
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		panic(err)
	}
	links := findLinks(doc, parsedUrl)
	fmt.Println(links)
}

func handleRelLinks(baseURL *url.URL, someURLStr string) string {
	parsedURL, err := baseURL.Parse(someURLStr)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return ""
	}
	// Strip fragments...
	parsedURL.Fragment = ""
	return parsedURL.String()
}

/*
We need to handle and parse relative links.
The process should be something like this
1) Check if the URL is an absolute URL
2) If so skip this URL (url.IsAbs)
3) Otherwise, we need to parse the relative URL using the base.
   Something like url.Parse(relativeUrl)
4) Write a tests to verify
*/

func findLinks(root *html.Node, baseURL *url.URL) []string {
	var links []string
	if root == nil {
		fmt.Println("nil node found!")
		return links
	}
	for node := range root.Descendants() {
		if node.Type == html.ElementNode && (node.DataAtom == atom.A || node.DataAtom == atom.Link) {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					links = append(links, handleRelLinks(baseURL, attr.Val))
				}
			}
		}
	}
	return links
}
