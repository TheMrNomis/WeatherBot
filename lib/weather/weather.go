package weather

import (
    "log"
    "bytes"
    "time"
    "os/exec"
    "net/http"
    "encoding/json"
    "strconv"
)

type OWM_WeatherResponse struct {
    Weather []OWM_Weather
    Base string
    Main OWM_Main
    Visibility int
    Wind OWM_Wind
    Clouds OWM_Cloud
    Dt int
    Sys OWM_Sys
    Id int
    Name string
    Cod int
}

type OWM_Weather struct {
    Id int
    Main string
    Description string
    Icon string
}

type OWM_Main struct {
    Temp float64
    Pressure float64
    Humidity float64
    Temp_min float64
    Temp_max float64
}

type OWM_Wind struct {
    Speed float64
    Deg float64
}

type OWM_Cloud struct {
    All int
}

type OWM_Sys struct {
    Type int
    Id int
    Message float64
    Country string
    Sunrise int
    Sunset int
}

func getWeather(city string) string {
    cmd := exec.Command("/bin/sh", "./weather.sh", city)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        log.Fatal(err)
    }
    return out.String()
}

var httpClient = &http.Client{Timeout: 10*time.Second}
func getJson(url string, target interface{}) error {
    r, err := httpClient.Get(url)
    defer r.Body.Close()
    if err != nil {
        return err
    }
    if r.StatusCode < 200 || r.StatusCode >= 300 {
        log.Println("Error while loading JSON:", r.StatusCode)
    }

    return json.NewDecoder(r.Body).Decode(target)
}

func GetWeatherResponse(city CityJson) (OWM_WeatherResponse, error) {
    url := "http://api.openweathermap.org/data/2.5/weather?id="+strconv.Itoa(city.Id)+"&lang=en&APPID=b2a51b191213f325c7709785ec4f6b72"
    log.Println(url)
    var weather OWM_WeatherResponse
    err := getJson(url, &weather)

    return weather, err
}

func Icon (weatherResponse OWM_WeatherResponse) string {
    var icon string
    switch weatherResponse.Weather[0].Icon {
    case "01d":
        icon = "☀️️"
    case "01n":
        icon = "🌕"
    case "02d":
        icon = "🌤️"
    case "02n":
        icon = "🌤️"
    case "03d":
        icon = "🌥️"
    case "03n":
        icon = "🌥️"
    case "04d":
        icon = "☁️️"
    case "04n":
        icon = "☁️️"
    case "09d":
        icon = "🌧️"
    case "09n":
        icon = "🌧️"
    case "10d":
        icon = "🌦️"
    case "10n":
        icon = "🌦️"
    case "11d":
        icon = "🌩️"
    case "11n":
        icon = "🌩️"
    case "13d":
        icon = "🌨️"
    case "13n":
        icon = "🌨️"
    case "50d":
        icon = "🌫️"
    case "50n":
        icon = "🌫️"
    }
}
