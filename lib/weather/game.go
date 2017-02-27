package weather

import (
    "log"
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

type GamePlayer struct {
    User *discordgo.User
    City CityJson
}

type GameGuild struct {
    Running bool
    Subject int
    Channel *discordgo.Channel
    Players []GamePlayer
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

        msg := "Game started on the subject of "

        switch currentGame.Subject {
        case Temp:
            msg += "Temperature"
        case Wind:
            msg += "Wind"
        case Rain:
            msg += "Rain"
        }

        session.ChannelMessageSend(currentGame.Channel.ID, msg)
        time.AfterFunc(30*time.Second, func() {
            endGame(currentGame, session, message)
        })
    }
}

func endGame(currentGame *GameGuild, session *discordgo.Session, message *discordgo.MessageCreate) {
    currentGame.Running = false
    currentGame.Players = []GamePlayer{}
    session.ChannelMessageSend(currentGame.Channel.ID, "game stopped")
    //TODO: tell the winner
    //TODO: remove all the propositions
}

func propose(messageArgs []string, currentGame *GameGuild, session *discordgo.Session, message *discordgo.MessageCreate) {
    if !currentGame.Running {
        session.ChannelMessageSend(currentGame.Channel.ID, "Cannot propose while game not started! (use `!wgstart` to start the game)")
        return
    }

    city, err := GetCityByArgs(m_db, messageArgs)
    if err != nil {
        log.Println("error retrieving city during !wgpropose: ", err)
        session.ChannelMessageSend(currentGame.Channel.ID, "Error, I could not understand the city :(")
        return
    }

    isNewProposition := true
    for i := 0; isNewProposition && i < len(currentGame.Players); i++ {
        if currentGame.Players[i].User.ID == message.Message.Author.ID {
            isNewProposition = false
            currentGame.Players[i].City = city
        }
    }

    if isNewProposition {
        tmpPlayer := GamePlayer{User: message.Message.Author, City: city}
        currentGame.Players = append(currentGame.Players, tmpPlayer)

        session.ChannelMessageSend(currentGame.Channel.ID, message.Message.Author.Username + " proposed " + countryNameToUT8Flag(city.Country) + " " + city.Name)
    } else {
        session.ChannelMessageSend(currentGame.Channel.ID, message.Message.Author.Username + " changed his mind for " + countryNameToUT8Flag(city.Country) + " " + city.Name)
    }
}
