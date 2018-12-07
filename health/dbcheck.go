package health

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DBDependency struct {
	BasicInfo DependencyInfo
	DSN       string
	Dialect   string
}

func dbStatusCheck(dbconfig *DBDependency) {
	sTime := time.Now()
	db, err := sql.Open(dbconfig.Dialect, dbconfig.DSN)
	if err != nil {
		dbconfig.BasicInfo.State.Status = CRIT
		dbconfig.BasicInfo.State.Details = fmt.Sprintf("%v", err)
	}
	_, err = db.Query(getSqlQuery("ping", dbconfig.Dialect))
	dbconfig.BasicInfo.ResponseTime = time.Since(sTime).Seconds()
	if err != nil {
		dbconfig.BasicInfo.State.Status = CRIT
		dbconfig.BasicInfo.State.Details = fmt.Sprintf("%v", err)
	} else {
		dbconfig.BasicInfo.State.Status = OK
		dbconfig.BasicInfo.State.Details = ""
	}
	dbconfig.BasicInfo.Version = fetchMySQLVersion(db, dbconfig.Dialect)
	dbconfig.BasicInfo.Type = "internal"
	defer db.Close()
}

func fetchMySQLVersion(db *sql.DB, dialect string) string {
	var version string
	rows, err := db.Query(getSqlQuery("version", dialect))
	if err != nil {
		return "Unknown"
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&version); err != nil {
			return "Unknown"
		}
	}
	return version
}

func getSqlQuery(sqltype string, dialect string) string {
	switch dialect {
	case "mysql":
		switch sqltype {
		case "version":
			return "SELECT version()"
		case "ping":
			return "SELECT 1"
		}
	case "postgre":
		switch sqltype {
		case "version":
			return "SELECT version()"
		case "ping":
			return "SELECT 1"
		}
	}
	return ""
}
