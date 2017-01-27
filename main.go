package main

import (
    "log"
    "fmt"
    "os"
    "os/signal"
    "flag"

    "github.com/bwmarrin/discordgo"
    "weatherbot/lib/botSettings"
    "weatherbot/lib/weather"
    "weatherbot/lib/install"
)

var m_settings botSettings.Settings
var m_install bool

func init() {
    settings, err := botSettings.GetSettings("settings.json")
    if err != nil {
        log.Fatal(err)
    }

    m_settings = settings
    weather.Init(&settings)

    flag.BoolVar(&m_install, "install", false, "--install if you want to install the db")

    flag.Parse()
}

func main() {
    if m_install {
        log.Println("installing")

        db, err := weather.OpenDatabase(&m_settings)
        if err != nil {
            log.Fatal(err)
        }
        install.CreateCityDb(db, "city.list.json")
    } else {
        defer weather.Close()
        if m_settings.Token == "" || m_settings.Token == "<your token>" {
            log.Fatal("no token provided")
            return
        }

        discord, err := discordgo.New(m_settings.Token)
        if err != nil {
            log.Fatal("Error creating Discord session: ", err)
            return
        }

        discord.AddHandler(weather.HandleReady)
        discord.AddHandler(weather.HandleGuildCreate)
        discord.AddHandler(weather.HandleChannelCreate)
        discord.AddHandler(weather.HandleMessage)

        err = discord.Open()
        if err != nil {
            log.Fatal("Error opening Discord session: ", err)
        }

        log.Println("WeatherBot is running, ^C to exit")
        sigchan := make(chan os.Signal, 10)
        signal.Notify(sigchan, os.Interrupt)
        <-sigchan
        fmt.Println()
        log.Println("Program killed")

        discord.Close()

        os.Exit(0)
    }
}
