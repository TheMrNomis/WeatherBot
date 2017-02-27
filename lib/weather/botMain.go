package weather

import (
    "log"
    "strings"
    "database/sql"

    "github.com/bwmarrin/discordgo"
    "weatherbot/lib/botSettings"
)

var m_settings *botSettings.Settings
var m_db *sql.DB

var m_thisUser *discordgo.User

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
        addGame(event.Channel)
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
            addGame(channel)
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

        if msgArgs[0] != "weather"{
            handleGameFunction(msgArgs, session, message)
        } else {
            city, err := GetCityByArgs(m_db, msgArgs)
            var msg string
            if err == nil {
                msg = GetWeatherStringForCity(city)
            } else {
                log.Println("error retrieving city: ", err)
                msg = "The city was not understood"
            }

            _, err = session.ChannelMessageSend(message.ChannelID, msg)
            if err != nil {
                log.Println("error sending message, ", err)
            }
        }
    }
}
