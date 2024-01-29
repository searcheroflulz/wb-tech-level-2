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

// сама функция wget
func wget(site string, root string) error {
	//создаем папку под названием сайта
	site = strings.TrimRight(site, "\r\n")
	siteURL, err := url.Parse(site)
	hostname := strings.TrimPrefix(siteURL.Hostname(), "www.")
	err = os.Mkdir(hostname, os.ModeDir)
	if err != nil {
		return err
	}
	//переходим в эту папку и ищем все ссылки на основной странице
	err = os.Chdir(filepath.Join(root, hostname))
	err = crawl(site, 2, site)

	return err
}

// рекурсивная функция для поиска ссылок и файлов на странице
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

	//проверяем, что ссылка не ведет на какой-то другой сайт и является сабдиректорией
	if isSameOrSubdirectory(baseURL, url) {
		err = downloadResources(resp.Body, baseURL)
		if err != nil {
			fmt.Println("Error downloading resources:", err)
		}
	}

	//достаем все ссылки со страницы
	links := extractLinks(downloadPage(url), baseURL)
	for _, link := range links {
		//преобразуем все ссылки в абсолютные ссылки
		absoluteURL, err := makeAbsoluteURL(link, baseURL)
		if err != nil {
			fmt.Println("Error making absolute URL:", err)
			continue
		}
		//продолжаем поиски ссылок на новых страницах
		err = crawl(absoluteURL, depth-1, baseURL)
		if err != nil {
			fmt.Println("Error crawling:", err)
		}
	}

	return nil
}

// функция для скачивания сурсов
func downloadResources(body io.Reader, baseURL string) error {
	tokenizer := html.NewTokenizer(body)

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			//обрабатываем теги, содержащие ресурсы, например, "img", "link", "script" и т.д.
			switch token.Data {
			case "img", "link", "script":
				for _, attr := range token.Attr {
					if attr.Key == "src" || attr.Key == "href" {
						resourceURL, err := makeAbsoluteURL(attr.Val, baseURL)
						if err != nil {
							fmt.Println("Error making absolute URL:", err)
							continue
						}
						err = downloadFile(resourceURL, baseURL)
						if err != nil {
							fmt.Println("Error downloading file:", err)
						}
					}
				}
			}
		}
	}
}

// скачиваем файлы по URL и сохранения в сабдиректории
func downloadFile(url string, baseURL string) error {
	if !isSameOrSubdirectory(baseURL, url) {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to fetch the resource. Status code: %d", resp.StatusCode)
	}

	//получаем имя файла
	fileName := filepath.Base(url)
	relPath, err := filepath.Rel(baseURL, url)
	if err != nil {
		return err
	}
	dirPath := filepath.Join(".", relPath)
	dirPath = filepath.Dir(dirPath)

	//создаем сабдиректорию, если ее нет
	if err := os.MkdirAll(dirPath, os.ModeDir); err != nil {
		return err
	}

	//создаем файл и сохраняем в нем данные
	file, err := os.Create(filepath.Join(dirPath, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Downloaded:", url)
	return nil
}

// преобразовываем URL в абсолютные
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

// ищем ссылки с помощью данной функции
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
						linkURL, err := makeAbsoluteURL(attr.Val, baseURL)
						if err == nil && isSameOrSubdirectory(baseURL, linkURL) {
							links = append(links, linkURL)
						}
					}
				}
			}
		}
	}
}

// проверяем, что сабдиректория относится к нашему основному сайту
func isSameOrSubdirectory(baseURL, linkURL string) bool {
	base, _ := url.Parse(baseURL)
	link, _ := url.Parse(linkURL)
	return strings.HasPrefix(link.Host, base.Host)
}

// скачиваем страницу
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
