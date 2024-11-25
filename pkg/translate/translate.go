package translate

import (
	"fmt"
	"github.com/abadojack/whatlanggo"
	"github.com/xiaoxuan6/deeplx"
	"strings"
	"time"
)

func Translation(description string) (string, bool) {
	if len(description) < 1 {
		return "", false
	}

	lang := whatlanggo.DetectLang(description)
	sourceLang := strings.ToUpper(lang.Iso6391())
	if strings.Compare(sourceLang, "zh") == 0 {
		return description, true
	}

	start := time.Now()
	response := deeplx.Translate(description, sourceLang, "zh")

	var data string
	if response.Code == 200 {
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
