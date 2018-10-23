package gin_middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"testing"
)

func TestDelimitedURLForStats(t *testing.T) {
	tables := []struct {
		original string
		idx      int
		expected string
	}{
		{"v1/entity/23", 2, "v1.entity"},
		{"/v1/entity/23", 2, "v1.entity"},
		{"/v1/entity/23", 1, "v1"},
		{"/v1/entity/23/value", 2, "v1.entity"},
	}

	for _, table := range tables {
		gin.SetMode(gin.TestMode)
		c := &gin.Context{}
		c.Request = &http.Request{}
		c.Request.URL = &url.URL{}
		c.Request.URL.Path = table.original

		c.Set(URL_TERMINATION_INDEX_VAR, table.idx)

		delimitedUrl := getDelimitedURLForStats(c)

		if delimitedUrl != table.expected {
			t.Errorf("Converted: %s, Expected: %s, Original: %s.", delimitedUrl, table.expected, table.original)
		}
	}
}
