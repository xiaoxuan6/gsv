package translate

import (
	"fmt"
	"github.com/xiaoxuan6/deeplx"
	"strings"
	"time"
)

func Translation(description string) string {
	start := time.Now()
	response := deeplx.Translate(description, "en", "zh")
	if response.Code != 200 {
		return description
	}
	end := time.Now().Sub(start).Seconds()

	return fmt.Sprintf(
		"%s {耗时：%s/s}",
		strings.TrimSpace(response.Data),
		fmt.Sprintf("%.2f", end),
	)
}
