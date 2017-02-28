package weather

import (
    "log"
    "time"
    "math/rand"
    "math"
    "github.com/bwmarrin/discordgo"
)

const (
    TempLow = iota
    TempHigh
    Wind
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
    EndTimer *time.Timer
    Players []GamePlayer
}

var m_currentGames []*GameGuild

func addGame(channel *discordgo.Channel) {
    tmp := GameGuild{Running: false, Channel: channel, EndTimer: nil}
    m_currentGames = append(m_currentGames, &tmp)
}

func computeScore(gameType int, weather OWM_WeatherResponse) float64 {
    switch gameType {
    case TempLow:
        return -weather.Main.Temp
    case TempHigh:
        return weather.Main.Temp
    case Wind:
        return weather.Wind.Speed
    }
    return 0.0
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

    if currentGame.EndTimer != nil {
        currentGame.EndTimer.Stop()
    }
    currentGame.EndTimer = time.AfterFunc(30*time.Second, func() {
        log.Println("stopping game")
        endGame(currentGame, session, message)
        currentGame.EndTimer = nil
    })

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
        case TempLow:
            msg += "Lowest temperature"
        case TempHigh:
            msg += "Highest temperature"
        case Wind:
            msg += "Wind"
        }

        session.ChannelMessageSend(currentGame.Channel.ID, msg)
    }
}

func enounceScores(currentGame *GameGuild, session *discordgo.Session) {
    if len(currentGame.Players) < 1 {
        session.ChannelMessageSend(currentGame.Channel.ID, "No one played :(")
        return
    }

    var maxScore float64
    var bestProposition *discordgo.User

    maxScore = -math.MaxFloat64
    for _, proposition := range currentGame.Players {
        weather, err := GetWeatherResponse(proposition.City)
        if err != nil {
            log.Println("error getting weather response: ", err)
        }
        score := computeScore(currentGame.Subject, weather)
        log.Println(proposition.User.Username, " scores ", score)

        if score >= maxScore {
            maxScore = score
            bestProposition = proposition.User

        }
    }

    session.ChannelMessageSend(currentGame.Channel.ID, "Winner is <@" + bestProposition.ID + ">")
}

func endGame(currentGame *GameGuild, session *discordgo.Session, message *discordgo.MessageCreate) {
    enounceScores(currentGame, session)
    currentGame.Running = false
    currentGame.Players = []GamePlayer{}
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

        session.ChannelMessageSend(currentGame.Channel.ID, "<@" + message.Message.Author.ID + "> proposed " + countryNameToUT8Flag(city.Country) + " " + city.Name)
    } else {
        session.ChannelMessageSend(currentGame.Channel.ID, "<@" + message.Message.Author.ID + "> changed his mind for " + countryNameToUT8Flag(city.Country) + " " + city.Name)
    }
}
