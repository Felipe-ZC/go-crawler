package main

import (
	"fmt"
	"net/http"
	"time"
	// "strings"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const ProtocolStr = "https://"

func main() {
	client := &http.Client{Timeout: 30 * time.Second}
	
	fmt.Println("Fetching...")
	resp, err := client.Get("https://github.com")
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
	links := findLinks(doc)
	fmt.Println(links)
}

func handleRelLinks(node *html.Node) string {
	return ""
}

func findLinks(root *html.Node) []string {
	var links []string
	if root == nil {
		fmt.Println("nil node found!")
		return links
	}
	for node := range root.Descendants() {
		if node.Type == html.ElementNode && (node.DataAtom == atom.A || node.DataAtom == atom.Link) {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
	}
	return links
}
