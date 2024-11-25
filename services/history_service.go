package services

import (
	"bufio"
	"io"
	"mvdan.cc/xurls/v2"
	"net/http"
	"strings"
	"time"
)

var (
	client = http.Client{
		Timeout: 3 * time.Second,
	}
	History = make([]string, 100)
	urls    = []string{
		"https://github-mirror.us.kg/https:/github.com/xiaoxuan6/go-package-example/blob/main/README_PHP.md",
		"https://github-mirror.us.kg/https:/github.com/xiaoxuan6/go-package-example/blob/main/README_OTHER.md",
		"https://github-mirror.us.kg/https:/github.com/xiaoxuan6/go-package-example/blob/main/README.md",
	}
)

func FetchHistory() {
	for _, url := range urls {
		wg.Add(1)

		url := url
		go func() {
			defer wg.Done()
			response, err := client.Get(url)
			if err != nil {
				return
			}

			defer response.Body.Close()
			f := bufio.NewReader(response.Body)
			for {
				line, _, err := f.ReadLine()
				if err == io.EOF {
					break
				}

				x := xurls.Relaxed()
				domain := x.FindString(string(line))
				domain = strings.ReplaceAll(domain, "github.com/", "")
				if len(domain) > 1 {
					History = append(History, domain)
				}
			}
		}()
	}

	wg.Wait()
}
