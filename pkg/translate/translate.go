package translate

import (
	"fmt"
	"github.com/OwO-Network/gdeeplx"
	"github.com/abadojack/whatlanggo"
	"strings"
	"time"
)

func Translation(description string) string {
	start := time.Now()
	info := whatlanggo.Detect(description)
	sourceLang := info.Lang.String()
	if sourceLang != "Mandarin" {
		response, err := gdeeplx.Translate(description, sourceLang, "zh", 0)
		if err != nil {
			//fmt.Println("translate err: ", err.Error())
			return description
		} else {
			end := time.Now().Sub(start).Seconds()
			description = fmt.Sprintf(
				"%s {耗时：%s/s}",
				strings.TrimSpace(response.(map[string]interface{})["data"].(string)),
				fmt.Sprintf("%.2f", end),
			)
		}
	}

	return description
}
