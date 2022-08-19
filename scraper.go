package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
	"sync"
)

func validateUrl(url string) bool {
	r, _ := regexp.Compile("^https?://.*")
	valid := r.MatchString(url)
	return valid
}

func fetchAllLinks(url string, baseDomain *string, baseUrls *[]string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := soup.Get(url)

	if err != nil {
		fmt.Println(err)

	} else {
		doc := soup.HTMLParse(resp)
		links := doc.FindAll("a")

		var newBaseUrls []string
		for _, link := range links {
			url := link.Attrs()["href"]
			if validateUrl(url) {
				//fmt.Println(link.Text(), " Link : ", url)
				if strings.Contains(url, *baseDomain) && !slices.Contains(*baseUrls, url) {
					*baseUrls = append(*baseUrls, url)
					newBaseUrls = append(newBaseUrls, url)
				}
			}
		}

		if newBaseUrls != nil {
			for _, url := range newBaseUrls {
				wg.Add(1)
				go fetchAllLinks(url, baseDomain, baseUrls, wg)
			}
		}

	}

}

func main() {
	fmt.Println("entering main")

	baseDomain := "go.dev"

	baseUrl := "https://" + baseDomain

	var baseUrls []string
	var wg sync.WaitGroup

	wg.Add(1)
	go fetchAllLinks(baseUrl, &baseDomain, &baseUrls, &wg)
	wg.Wait()

	for _, link := range baseUrls {
		wg.Add(1)
		go fetchAllLinks(link, &baseDomain, &baseUrls, &wg)
	}

	wg.Wait()

	for _, link := range baseUrls {
		fmt.Println(link)
	}

}
