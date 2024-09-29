package translate

import (
	"bytes"
	"fmt"
	"github.com/OwO-Network/gdeeplx"
	"github.com/tidwall/gjson"
	"github.com/xiaoxuan6/deeplx"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func Translation(description string) (string, bool) {
	start := time.Now()
	response := deeplx.Translate(description, "en", "zh")

	var data string
	if response.Code != 200 {
		rand.Seed(time.Now().UnixNano())
		randomBit := rand.Intn(2)
		switch randomBit {
		case 0:
			data = translate(description)
		case 1:
			data = missuoTranslate(description)
		default:
			data = description
		}
	} else {
		data = response.Data
	}

	if len(data) < 1 {
		return description, false
	}

	end := time.Now().Sub(start).Seconds()
	return fmt.Sprintf(
		"%s {耗时：%s/s}",
		strings.TrimSpace(data),
		fmt.Sprintf("%.2f", end),
	), true
}

func translate(description string) string {
	result, err := gdeeplx.Translate(description, "EN", "ZH", 0)
	if err != nil {
		return description
	}

	return result.(map[string]interface{})["data"].(string)
}

// missuoTranslate 需要设置 key 并且使用 proxy 开启 7890 端口
func missuoTranslate(description string) string {
	key := os.Getenv("TRANSLATE_KEY")
	if key == "" {
		return description
	}

	body := bytes.NewBufferString(`{"text":"` + description + `","source_lang":"en","target_lang":"zh"}`)

	proxyUrl, _ := url.Parse("http://127.0.0.1:7890")
	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	response, err := client.Post(fmt.Sprintf("https://deeplx.missuo.ru/translate?key=%s", key), "application/json", body)
	if err != nil {
		return description
	}
	defer response.Body.Close()

	b, _ := ioutil.ReadAll(response.Body)
	if code := gjson.ParseBytes(b).Get("code").Int(); code == 200 {
		return gjson.ParseBytes(b).Get("data").String()
	}

	return description
}
