package logging

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

//You must call InitStandardLogger at the beginning of your application
//to set up logging with the standard format.

func InitStandardLogger(version string) {
	InitLogger(version, nil)
}

func InitLogger(version string, l *logrus.Logger) {
	if l == nil {
		l = logrus.StandardLogger()
	}
	l.SetFormatter(&CustomJSONFormatter{version: version})
}

type CustomJSONFormatter struct {
	version string
}

func (f *CustomJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Data["timestamp"] = entry.Time.UTC().Format("2006-01-02 15:04:05 -0700")
	entry.Data["level"] = entry.Level.String()
	entry.Data["version"] = f.version

	if entry.Level <= logrus.ErrorLevel {
		entry.Data["message"] = entry.Message
	} else if entry.Message != "" {
		entry.Data["message"] = entry.Message
	}

	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON: %v", err)
	}
	return append(serialized, '\n'), nil
}

//RL (request logging) records correlation id. You must use middleware.Correlation()
//as a middleware in order for the correlation ID to appear.
func RL(h *http.Request, w http.ResponseWriter) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"correlation_id": w.Header().Get("X-CORRELATION-ID"),
	})
}
