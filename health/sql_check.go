package health

import (
	"database/sql"
	"errors"
	"time"
)

type SQLCheck struct {
	Name string
	Type string
	DB   interface{}
}

func (sc SQLCheck) Check() *DependencyInfo {
	var err error
	var t float64
	var version string
	state := DependencyState{Status: OK}

	switch conn := sc.DB.(type) {
	case *sql.DB:
		sTime := time.Now()
		row := conn.QueryRow("SELECT version()")
		t = time.Since(sTime).Seconds()

		if err := row.Scan(&version); err != nil {
			state.Status = WARN
			state.Details = "Error retrieving version: " + err.Error()
		}
	default:
		err = errors.New("Unknown type of DB connection")
	}

	if err != nil {
		state.Status = CRIT
		state.Details = err.Error()
	}
	return &DependencyInfo{
		Name:         sc.Name,
		Type:         "internal",
		Version:      version,
		Revision:     "",
		State:        state,
		ResponseTime: t,
	}
}

func (sc SQLCheck) GetName() string {
	return sc.Name
}

func (sc SQLCheck) GetType() string {
	return sc.Type
}
