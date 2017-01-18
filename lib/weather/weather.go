package weather

import (
    "log"
    "regexp"
    "bytes"
    "os/exec"

    "github.com/bwmarrin/discordgo"
    "weatherbot/lib/botSettings"
)

var m_settings *botSettings.Settings

var m_thisUser *discordgo.User
var m_weatherChannels []*discordgo.Channel

func getWeather(city string) string {
    cmd := exec.Command("/bin/sh", "./weather.sh", city)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        log.Fatal(err)
    }
    return out.String()
}


func Init(settings *botSettings.Settings) {
    m_settings = settings;
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
    var rgxp = regexp.MustCompile(`^!weather(\s(\w*))?$`)
    if rgxp.MatchString(message.Content) {
        result := rgxp.FindStringSubmatch(message.Content)
        var city = result[2]

        if city == "" {
            city = m_settings.DefaultLocation
        }

        log.Println("city = '", city, "'")

        var msg = "```" + getWeather(city) + "```"
        _, err := session.ChannelMessageSend(message.ChannelID, msg)
        if err != nil {
            log.Println("error sending message, ", err)
        }
    }
}
