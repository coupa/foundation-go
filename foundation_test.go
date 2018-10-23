package foundation

import (
	"github.com/coupa/foundation-go/logging"
	"testing"
)

func TestInit(t *testing.T) {
	tables := []struct {
		project string
		app     string
		ver     string
	}{
		{"p", "a", "v"},
	}

	for _, table := range tables {
		InitMetricsMonitoring(table.project, table.app, table.ver)
		if logging.LoggingApp != table.app {
			t.Errorf("Logging app not initialized")
		}
		if logging.LoggingAppVersion != table.ver {
			t.Errorf("Logging ver not initialized")
		}
		if logging.LoggingProject != table.project {
			t.Errorf("Logging project not initialized")
		}
	}
}
