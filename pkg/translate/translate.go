package translate

import (
	"github.com/OwO-Network/gdeeplx"
	"github.com/abadojack/whatlanggo"
	"strings"
)

func Translation(description string) string {
	info := whatlanggo.Detect(description)
	sourceLang := info.Lang.String()
	if sourceLang != "Mandarin" {
		response, err := gdeeplx.Translate(description, sourceLang, "zh", 0)
		if err != nil {
			//fmt.Println("translate err: ", err.Error())
			return description
		} else {
			description = strings.TrimSpace(response.(map[string]interface{})["data"].(string))
		}
	}

	return description
}
