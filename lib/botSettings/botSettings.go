package botSettings

import (
    "encoding/json"
    "errors"
    "log"
    "os"
)

type Settings struct {
    Token string
    APIkey string
    Channel string
    DefaultLocation string
}

func setError(err error) error {
    return errors.New("BotSetting: " + err.Error())
}

func GetSettings(filename string) (Settings, error) {
    file, err := os.Open(filename)
    if err != nil {
        localError := setError(err)
        log.Fatal(localError)
        return Settings{}, localError
    }
    dec := json.NewDecoder(file)

    var settings Settings
    if err := dec.Decode(&settings); err != nil {
        localError := setError(err)
        log.Fatal(localError)
        return Settings{}, localError
    }

    return settings, nil
}
