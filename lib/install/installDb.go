package install

import (
    "log"
    "os"

    "encoding/json"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"

    "weatherbot/lib/weather"
)

func CreateCityDb (db *sql.DB, cityFilename string) error {
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS cities (city_id integer not null primary key, city_name text, city_country text);`)
    if err != nil {
        log.Fatal("Error creating table city:", err)
        return err
    }

    cityFile, err := os.Open(cityFilename)
    defer cityFile.Close()
    if err != nil {
        log.Fatal("Error loading the city file:", err)
        return err
    }

    dec := json.NewDecoder(cityFile)
    var cities []weather.CityJson
    if err := dec.Decode(&cities); err != nil {
        log.Fatal("Error decoding the json:", err)
        return err
    }

    log.Println("Opened", len(cities), "cities")

    tx, err := db.Begin()
    if err != nil {
        log.Fatal("Error beginning transaction:", err)
        return err
    }

    stmt, err := tx.Prepare("INSERT OR IGNORE INTO cities(city_id, city_name, city_country) VALUES(?,?,?);")
    defer stmt.Close()
    if err != nil {
        log.Fatal("Error preparing statement:", err)
        return err
    }

    for _,city := range cities {
        _, err = stmt.Exec(city.Id, city.Name, city.Country)
        if err != nil {
            log.Fatal("Error writing city", city.Name, ":", err)
        }
    }
    tx.Commit()

    return nil
}
