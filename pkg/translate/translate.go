package translate

import (
	"fmt"
	"github.com/xiaoxuan6/deeplx"
	"strings"
	"time"
)

func Translation(description string) (string, bool) {
	start := time.Now()
	response := deeplx.Translate(description, "en", "zh")
	if response.Code != 200 {
		return description, false
	}

	if len(response.Data) < 1 {
		return description, false
	}

	end := time.Now().Sub(start).Seconds()
	return fmt.Sprintf(
		"%s {耗时：%s/s}",
		strings.TrimSpace(response.Data),
		fmt.Sprintf("%.2f", end),
	), true
}
