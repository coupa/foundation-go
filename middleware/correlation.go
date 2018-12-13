package middleware

import (
	// log "github.com/coupa/foundation-go/logging"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
)

const correlationHeader = "X-CORRELATION-ID"

func Correlation() gin.HandlerFunc {
	return func(c *gin.Context) {
		SetCorrelation(c.Request, c.Writer)
		c.Next()
	}
}

func SetCorrelation(req *http.Request, resp http.ResponseWriter) {
	if req.URL.Path == "/health" {
		return
	}
	correlationID := req.Header.Get(correlationHeader)
	if correlationID == "" {
		//This will change on go.uuid master, where uuid.NewV4() returns the id and an error
		correlationID = uuid.NewV4().String()
	}
	if correlationID != "" {
		resp.Header().Set(correlationHeader, correlationID)
	}
}
