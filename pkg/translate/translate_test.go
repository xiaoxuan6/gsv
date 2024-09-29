package translate

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMissuoTranslate(t *testing.T) {
	_ = os.Setenv("TRANSLATE_KEY", "xxx")
	result := missuoTranslate("hello word")
	assert.Equal(t, result, "你好")
	t.Log(result)
	assert.Nil(t, nil)
}
