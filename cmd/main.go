package main

import (
	"LinkScrapper/graph"
	"LinkScrapper/scrapper"
	"fmt"
)

func main() {
	var link string

	graph.SetGraphDimensions(5, 5)

	fmt.Println("Enter the wiki link")
	fmt.Scan(&link)

	root := scrapper.Scrap(link)

	root.GraphToJson("GraphicalLinks.json")
}
