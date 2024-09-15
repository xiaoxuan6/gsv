package deeplx

import (
	"bufio"
	"encoding/json"
	"github.com/abadojack/whatlanggo"
	"github.com/avast/retry-go"
	"github.com/tidwall/gjson"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	targetUrls = make([]string, 0)
	urls       = []string{"https://deeplx.mingming.dev/translate"}
)

type request struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

type Response struct {
	Code int64  `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

func fetchUri() string {
	if len(targetUrls) < 1 {
		client := &http.Client{
			Timeout: 3 * time.Second,
		}

		resp, err := client.Get("https://github-mirror.us.kg/https://github.com/ycvk/deeplx-local/blob/windows/url.txt")
		if err == nil && resp.StatusCode == 200 {
			r := bufio.NewReader(resp.Body)
			for {
				line, _, errs := r.ReadLine()
				if errs == io.EOF {
					break
				}

				targetUrls = append(targetUrls, string(line))
			}
			urls = append(urls, targetUrls...)
		}
	}

	urlsLen := len(urls)
	randomIndex := rand.Intn(urlsLen)
	if randomIndex >= urlsLen {
		return urls[0]
	} else {
		return urls[randomIndex]
	}
}

func Translate(text, sourceLang, targetLang string) *Response {
	if len(text) == 0 {
		return &Response{
			Code: 500,
			Msg:  "No Translate Text Found",
		}
	}

	if len(sourceLang) == 0 {
		lang := whatlanggo.DetectLang(text)
		sourceLang = strings.ToUpper(lang.Iso6391())
	}

	if len(targetLang) == 0 {
		targetLang = "EN"
	}

	req := &request{
		Text:       text,
		SourceLang: sourceLang,
		TargetLang: targetLang,
	}
	jsonBody, _ := json.Marshal(req)

	var body []byte
	_ = retry.Do(
		func() error {
			response, err := http.Post(fetchUri(), "application/json", strings.NewReader(string(jsonBody)))

			if err == nil {
				defer func() {
					_ = response.Body.Close()
				}()

				body, err = io.ReadAll(response.Body)
			} else {
				body = []byte(`{"code":500, "message": ` + err.Error() + `}`)
			}

			return err
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	)

	return &Response{
		Code: gjson.Get(string(body), "code").Int(),
		Data: gjson.Get(string(body), "data").String(),
		Msg:  gjson.Get(string(body), "message").String(),
	}
}