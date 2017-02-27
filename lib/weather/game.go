package weather

import (
    "time"
    "math/rand"
    "github.com/bwmarrin/discordgo"
)

const (
    Temp = iota
    Wind
    Rain
    EndWeatherSubjectEnum
)

type GameGuild struct {
    Running bool
    Subject int
    Channel *discordgo.Channel
}

var m_currentGames []*GameGuild

func addGame(channel *discordgo.Channel) {
    tmp := GameGuild{Running: false, Channel: channel}
    m_currentGames = append(m_currentGames, &tmp)
}

func handleGameFunction(messageArgs []string, session *discordgo.Session, message *discordgo.MessageCreate) {
    var currentGame *GameGuild
    currentGameFound := false
    for _, game := range m_currentGames {
        if message.ChannelID == game.Channel.ID {
            currentGame = game
            currentGameFound = true
            break
        }
    }

    if !currentGameFound {
        session.ChannelMessageSend(message.ChannelID, "You're not in the right channel!")
        return
    }

    switch messageArgs[0] {
    case "wgstart":
        startGame(currentGame, session, message)
    case "wgpropose":
        propose(messageArgs, currentGame, session, message)
    }
}

func startGame(currentGame *GameGuild, session *discordgo.Session, message *discordgo.MessageCreate) {
    rand.Seed(time.Now().UnixNano())
    if currentGame.Running {
        session.ChannelMessageSend(currentGame.Channel.ID, "Game is already running.")
    } else {
        currentGame.Running = true
        currentGame.Subject = rand.Intn(EndWeatherSubjectEnum)

        message := "Game started on the subject of "

        switch currentGame.Subject {
        case Temp:
            message += "Temperature"
        case Wind:
            message += "Wind"
        case Rain:
            message += "Rain"
        }

        session.ChannelMessageSend(currentGame.Channel.ID, message)
        time.AfterFunc(10*time.Second, func() {
            currentGame.Running = false
            session.ChannelMessageSend(currentGame.Channel.ID, "game stopped")
            //TODO: tell the winner
        })
    }
}

func propose(messageArgs []string, currentGame *GameGuild, session *discordgo.Session, message *discordgo.MessageCreate) {
}
