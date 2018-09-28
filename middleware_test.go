package common_go

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"testing"
)

func TestToIndexName(t *testing.T) {
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

		c.Set(ID_INDEX_IN_URL_VAR, table.idx)

		url := getDelimitedURLForStats(c)

		if url != table.expected {
			t.Errorf("Converted: %d, Expected: %d, Original: %d.", url, table.expected, table.original)
		}
	}
}
