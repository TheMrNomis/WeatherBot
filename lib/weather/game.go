package weather

import (
    "github.com/bwmarrin/discordgo"
)

type WeatherSubject int
const (
    Temp WeatherSubject = iota
    Wind
    Rain
)

type GameGuild struct {
    Running bool
    Subject WeatherSubject
    Guild *discordgo.Guild
}
