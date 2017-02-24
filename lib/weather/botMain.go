package weather

import (
    "log"
    "strings"
    "database/sql"
    "unicode/utf8"

    "github.com/bwmarrin/discordgo"
    "weatherbot/lib/botSettings"
)

var m_settings *botSettings.Settings
var m_db *sql.DB

var m_thisUser *discordgo.User
var m_weatherChannels []*discordgo.Channel

func Init(settings *botSettings.Settings) {
    m_settings = settings;

    db, err := OpenDatabase(settings)
    m_db = db
    if err != nil {
        log.Fatal(err)
    }
}

func Close() {
    m_db.Close()
}

func HandleReady(session *discordgo.Session, event *discordgo.Ready) {
    m_thisUser = event.User
}

func HandleChannelCreate(session *discordgo.Session, event *discordgo.ChannelCreate) {
    if event.Channel.Name == m_settings.Channel {
        m_weatherChannels = append(m_weatherChannels, event.Channel)
    }
}

func HandleGuildCreate(session *discordgo.Session, event *discordgo.GuildCreate) {
    if (*event).Guild.Unavailable {
        log.Println("unavailable guild: ", event.Guild.Name)
        return
    }

    log.Println("joined guild", event.Guild.Name)
    var channelFound bool = false;
    for _, channel := range event.Guild.Channels {
        if !channelFound && channel.Name == m_settings.Channel {
            m_weatherChannels = append(m_weatherChannels, channel)
            channelFound = true
        }
    }
    if !channelFound {
        log.Println("no channel '", m_settings.Channel, "' found")
    }
}

func HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
    trimedMessage := strings.TrimPrefix(message.Content, "!")
    if trimedMessage != message.Content { //if they are equals, the prefix was not there
        msgArgs := strings.Split(strings.Replace(trimedMessage, ",", " ", -1), " ")

        switch msgArgs[0] {
        case "weather":
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
            var msg string
            if country != "" {
                msg = getWeatherWithCountry(city, country)
            } else {
                if city == "" {
                    city = m_settings.DefaultLocation
                }
                msg = getWeather(city)
            }

            _, err := session.ChannelMessageSend(message.ChannelID, msg)
            if err != nil {
                log.Println("error sending message, ", err)
            }
        }
    }
}
