package weather

import (
    "strings"
    "unicode/utf8"

    "database/sql"
    _ "github.com/mattn/go-sqlite3"

    "weatherbot/lib/botSettings"
)

type CityJson struct {
    Id int
    Name string
    Country string
}

func OpenDatabase(settings *botSettings.Settings) (*sql.DB, error) {
    return sql.Open("sqlite3", settings.DBFile)
}

func GetCityByArgs(db *sql.DB, msgArgs []string) (CityJson, error) {
    city := "";
    country := "";

    for _, str := range msgArgs[1:len(msgArgs)] {
        strlen := utf8.RuneCountInString(str)
        if strlen > 0 {
            if city == "" || strlen > 2 {
                city = str
            } else {
                country = str
            }
        }
    }
    if country != "" {
        return GetCityByNameAndCountry(db, city, country)
    } else {
        if city == "" {
            city = m_settings.DefaultLocation
        }
        return GetCityByName(db, city)
    }
}

func GetCityByName(db *sql.DB, cityName string) (CityJson, error) {
    cityName = formatCityName(cityName)

    var cityId int
    var cityCountry string
    err := db.QueryRow(`SELECT city_id, city_country FROM cities WHERE city_name = ? ORDER BY city_id LIMIT 0,1;`, cityName).Scan(&cityId, &cityCountry)

    return CityJson{cityId, cityName, cityCountry}, err
}

func GetCityByNameAndCountry(db *sql.DB, cityName string, cityCountry string) (CityJson, error) {
    cityName = formatCityName(cityName)
    cityCountry = formatCountryName(cityCountry)

    var cityId int
    err := db.QueryRow(`SELECT city_id FROM cities WHERE city_name = ? AND city_country = ? ORDER BY city_id LIMIT 0,1;`, cityName, cityCountry).Scan(&cityId)

    return CityJson{cityId, cityName, cityCountry}, err
}

func formatCityName(cityName string) string {
    return strings.Title(cityName)
}

func formatCountryName(countryName string) string {
    countryName = strings.ToUpper(countryName)

    switch countryName {
    case "UK":
        countryName = "GB"
    }

    return countryName;
}
