package crawler

import (
	"sync"

	"github.com/gocolly/colly/v2"
)

type Crawler struct {
	Collector *colly.Collector
	Mutex     sync.Mutex
}

func NewCrawler() *Crawler {
	c := colly.NewCollector(
		colly.AllowedDomains("pokedex.org"),
	)
	c.UserAgent = "Colly - Golang scraping framework"
	return &Crawler{
		Collector: c,
	}
}

func (cr *Crawler) Start(urls []string, handleFunc func(*colly.Collector)) {
	for _, url := range urls {
		// Clone the collector to avoid reusing handlers
		cloned := cr.Collector.Clone()
		handleFunc(cloned)
		cloned.Visit(url)
	}
}
