package foundation

import (
	"testing"
	"github.com/coupa/foundation-go/logging"
)

func TestFoundation(t *testing.T) {
}

func TestInit(t *testing.T) {
	tables := []struct {
		project string
		app   string
		ver  string
		component string
	}{
		{"p", "a", "v", "c"},
	}

	for _, table := range tables {
		InitMetricsMonitoringOnComponent(table.project, table.app, table.ver, table.component)
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