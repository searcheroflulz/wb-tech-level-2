package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func wget(site string, root string) error {
	site = strings.TrimRight(site, "\r\n")
	siteURL, err := url.Parse(site)
	hostname := strings.TrimPrefix(siteURL.Hostname(), "www.")
	err = os.Mkdir(hostname, os.ModeDir)
	if err != nil {
		return err
	}
	err = os.Chdir(filepath.Join(root, hostname))
	err = crawl(site, 2, site)

	return err
}

func crawl(url string, depth int, baseURL string) error {
	if depth <= 0 {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to fetch the page %s. Status code: %d", url, resp.StatusCode)
	}

	links := extractLinks(downloadPage(url), baseURL)
	for _, link := range links {
		// Преобразование относительных URL в абсолютные
		absoluteURL, err := makeAbsoluteURL(link, baseURL)
		if err != nil {
			fmt.Println("Error making absolute URL:", err)
			continue
		}

		err = crawl(absoluteURL, depth-1, baseURL)
		if err != nil {
			fmt.Println("Error crawling:", err)
		}
	}

	return nil
}

// Пример простой функции для преобразования относительных URL в абсолютные
func makeAbsoluteURL(relativeURL, baseURL string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	rel, err := url.Parse(relativeURL)
	if err != nil {
		return "", err
	}

	return base.ResolveReference(rel).String(), nil
}

func extractLinks(body []byte, baseURL string) []string {
	var links []string
	bodyReader := bytes.NewReader(body)
	z := html.NewTokenizer(bodyReader)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						linkURL, err := resolveURL(attr.Val, baseURL)
						if err == nil && isSameOrSubdirectory(baseURL, linkURL) {
							links = append(links, linkURL)
						}
					}
				}
			}
		}
	}
}

func resolveURL(href, baseURL string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	rel, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	absURL := base.ResolveReference(rel)
	return absURL.String(), nil
}

func isSameOrSubdirectory(baseURL, linkURL string) bool {
	base, _ := url.Parse(baseURL)
	link, _ := url.Parse(linkURL)
	return strings.HasPrefix(link.Host, base.Host)
}

func downloadPage(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	filename := filepath.Base(url)
	file, err := os.Create(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	_, err = file.Write(data)
	if err != nil {
		return nil
	}

	return data
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	site, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = wget(site, path)
	if err != nil {
		log.Fatal(err)
	}

}
