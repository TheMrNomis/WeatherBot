package weather

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"

    "weatherbot/lib/botSettings"
)

func OpenDatabase(settings botSettings.Settings) sql.DB, err {
    return sql.Open("sqlite3", settings.DBFile)
}
