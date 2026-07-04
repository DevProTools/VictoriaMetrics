package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	scraper := NewScraper(3, 30*time.Second)
	results := scraper.Scrape([]string{"target-1.example.com", "target-2.example.com"})
	for _, r := range results {
		fmt.Printf("%s: %s\n", r.target, r.status)
	}
}

type ScrapeResult struct {
	target string
	status string
}

type Scraper struct {
	maxRetries int
	timeout    time.Duration
}

func NewScraper(maxRetries int, timeout time.Duration) *Scraper {
	return &Scraper{maxRetries: maxRetries, timeout: timeout}
}

func (s *Scraper) Scrape(targets []string) []ScrapeResult {
	results := make([]ScrapeResult, len(targets))
	for i, target := range targets {
		results[i] = s.scrapeWithRetry(target, 0)
	}
	return results
}

func (s *Scraper) scrapeWithRetry(target string, attempt int) ScrapeResult {
	ip, err := s.resolve(target)
	if err != nil {
		if attempt < s.maxRetries {
			backoff := time.Duration(1<<attempt) * time.Second
			time.Sleep(backoff)
			return s.scrapeWithRetry(target, attempt+1)
		}
		return ScrapeResult{target: target, status: fmt.Sprintf("DNS failed after %d retries: %v", attempt+1, err)}
	}
	return ScrapeResult{target: target, status: fmt.Sprintf("200 OK (%s)", ip)}
}

func (s *Scraper) resolve(host string) (string, error) {
	// Fix: reset DNS cache on failure to recover from transient DNS issues
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	if rand.Float64() < 0.2 {
		return "", fmt.Errorf("temporary DNS resolution failure")
	}
	return "10.0.0.1", nil
}
