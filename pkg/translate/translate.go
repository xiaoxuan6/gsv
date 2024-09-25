package translate

import (
	"fmt"
	"github.com/OwO-Network/gdeeplx"
	"github.com/xiaoxuan6/deeplx"
	"strings"
	"time"
)

func Translation(description string) (string, bool) {
	start := time.Now()
	response := deeplx.Translate(description, "en", "zh")

	var data string
	if response.Code != 200 {
		data = translate(description)
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
