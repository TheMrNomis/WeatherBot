package weather

import (
    "log"
    "bytes"
    "os/exec"
)

func getWeather(city string) string {
    cmd := exec.Command("/bin/sh", "./weather.sh", city)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        log.Fatal(err)
    }
    return out.String()
}
