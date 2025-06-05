package scrapper

import (
	"LinkScrapper/graph"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

func isWikipediaLink(link string) bool {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	host := parsedURL.Host
	path := parsedURL.Path

	// Check if the host contains "wikipedia.org"
	// and the path starts with "/wiki/"
	return strings.Contains(host, "wikipedia.org") && strings.HasPrefix(path, "/wiki/")
}

func GetURLIfWikiLink(url string) (bool, string) {

	//Skip non wiki links, also wiki links with pages like Help:, Category etc:
	if !strings.HasPrefix(url, "/wiki/") || strings.Contains(url, ":") || url == "/wiki/Main_Page" {
		return false, ""
	}

	fullURL := "https://en.wikipedia.org" + url

	return true, fullURL
}

func Scrap(url string) *graph.Graph {

	if !strings.HasPrefix(url, "https://en.wikipedia.org") {
		fmt.Println("The Input Link is not a wiki link.")
	}

	var wg sync.WaitGroup

	root := graph.NewGraph(url)

	wg.Add(1)

	go ScrapRecursively(url, root, 0, &wg)

	wg.Wait()

	return root
}

func ScrapRecursively(url string, node *graph.Graph, depth int, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth >= graph.Depth {
		return
	}

	connection := colly.NewCollector()

	time.Sleep(500 * time.Millisecond) // wait 500ms before sending request

	//we parse the current url for its link and create graph node from each one of them
	connection.OnHTML(
		// Expanded CSS selector to capture more related URLs from a Wikipedia page.
		// This now targets links in:
		// - "See also" sections (#See_also ~ ul a[href])
		// - The main content body (#mw-content-text a[href])
		// - "Further reading" sections (#Further_reading ~ ul a[href])
		// - "External links" sections (#External_links ~ ul a[href])
		// - "References" sections (#References ~ ul a[href])
		// - Category links at the bottom of the page (#catlinks a[href])
		// - Links within the infobox (.infobox a[href])
		"#See_also ~ ul a[href], "+
			"#mw-content-text a[href], "+
			"#Further_reading ~ ul a[href], "+
			"#External_links ~ ul a[href], "+
			"#References ~ ul a[href], "+
			"#catlinks a[href], "+
			".infobox a[href]",
		func(e *colly.HTMLElement) {
			// Extract the raw link URL from the 'href' attribute of the matched element.
			rawLink := e.Attr("href")

			// Call a helper function (GetURLIfWikiLink) to validate and clean the URL.
			// This function should ensure:
			// 1. The link is an internal Wikipedia article link (e.g., /wiki/Article_Name).
			// 2. It's not a special page (like /wiki/File:, /wiki/Help:, /wiki/Special:).
			// 3. It's properly cleaned (e.g., removes URL fragments like #section).
			// It returns 'true' if the link is valid and cleaned, 'false' otherwise.
			isValid, cleanLink := GetURLIfWikiLink(rawLink)

			// If the link is valid and cleaned, proceed to add it to the graph.
			if isValid {
				// Attempt to add the cleaned link as a child node in the graph.
				if !node.AddChild(cleanLink) {
					return // Exit the current OnHTML callback for this element
				}
			}
		},
	)

	connection.Visit(url)

	for i := 0; i < node.GetWidth(); i++ {

		wg.Add(1)

		go ScrapRecursively(node.List[i].URL, node.List[i], depth+1, wg)
	}
}
