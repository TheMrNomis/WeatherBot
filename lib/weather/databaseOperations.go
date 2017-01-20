package weather

import (
    //"log"

    "database/sql"
    _ "github.com/mattn/go-sqlite3"

    "weatherbot/lib/botSettings"
)

type CityJson struct {
    Id int
    Name string
    Country string
}

func OpenDatabase(settings botSettings.Settings) (*sql.DB, error) {
    return sql.Open("sqlite3", settings.DBFile)
}

func GetCityByName(db *sql.DB, cityName string) (CityJson, error) {
    var cityId int
    var cityCountry string
    err := db.QueryRow(`SELECT city_id, city_country FROM cities WHERE city_name = ? ORDER BY city_id LIMIT 0,1;`, cityName).Scan(&cityId, &cityCountry)

    return CityJson{cityId, cityName, cityCountry}, err
}

func GetCityByNameAndCountry(db *sql.DB, cityName string, cityCountry string) (CityJson, error) {
    var cityId int
    err := db.QueryRow(`SELECT city_id FROM cities WHERE city_name = ? AND city_country = ? ORDER BY city_id LIMIT 0,1;`, cityName, cityCountry).Scan(&cityId)

    return CityJson{cityId, cityName, cityCountry}, err
}
