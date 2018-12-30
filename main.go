package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Variables used for command line parameters
var (
	Token string
)

type Rankings struct {
	Admin             bool     `json:"admin"`
	BannerUrl         string   `json:"banner_url"`
	Guild             Guild    `json:"guild"`
	Page              int      `json:"page"`
	Player            Player   `json:"player"`
	Players           []Player `json:"players"`
	RoleRewards       []string `json:"role_rewards"`
	UserGuildSettings string   `json:"user_guild_settings"`
}

type Guild struct {
	Icon    string `json:"icon"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Premium bool   `json:"premium"`
}

type Player struct {
	Avatar        string `json:"avatar"`
	DetailedXp    []int  `json:"detailed_xp"`
	Discriminator string `json:"discriminator"`
	GuildId       string `json:"guild_id"`
	Id            string `json:"id"`
	Level         int    `json:"level"`
	Username      string `json:"username"`
	Xp            int64  `json:"xp"`
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New(Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	//go startFarming(dg)

	// Register the messageCreate func as a callback for MessageCreate events.
	// TODO: add debounce
	dg.AddHandler(handleRanking)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	delay := time.Duration((rand.Intn(5) + 1) * 5)

	// If the message is "ping" reply with "Pong!"
	if strings.Contains(m.Content, ":qstory:") {
		<-time.After(delay * time.Second)
		s.ChannelMessageSend(m.ChannelID, "<:qstory:527743483163967498>")
	}

	// If the message is "pong" reply with "Ping!"
	if strings.Contains(m.Content, ":qhi:") {
		<-time.After(delay * time.Second)
		s.ChannelMessageSend(m.ChannelID, "<:qhi:527743420689678356>")
	}

	// If the message is "pong" reply with "Ping!"
	if strings.Contains(m.Content, ":qlul:") {
		<-time.After(delay * time.Second)
		s.ChannelMessageSend(m.ChannelID, "<:qlul:527743432014168064>")
	}
}

func startFarming(s *discordgo.Session) {
	for {
		msg, _ := s.ChannelMessageSend("527751887185903626", "просто фармим опыт...")
		<-time.After(2 * time.Minute)
		s.ChannelMessageDelete("527751887185903626", msg.ID)
	}
}

func handleRanking(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!ranking" {
		response, err := http.Get("https://mee6.xyz/api/plugins/levels/leaderboard/404185377264369665")
		if err != nil {
			fmt.Println("error retrieving dat,", err)
			return
		}

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)

		var rankings Rankings
		err = json.Unmarshal(contents, &rankings)

		if err != nil {
			fmt.Println("error parsing json,", err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"TOP:\n1. %s (%d level with %d messages)\n2. %s (%d level with %d messages)\n3. %s (%d level with %d messages)",
			rankings.Players[0].Username, rankings.Players[0].Level, rankings.Players[0].Xp/20,
			rankings.Players[1].Username, rankings.Players[1].Level, rankings.Players[1].Xp/20,
			rankings.Players[2].Username, rankings.Players[2].Level, rankings.Players[2].Xp/20))
	}

}

//func startMusic(s *discordgo.Session) {
//	for {
//		<-time.After(2 * time.Minute)
//		go s.ChannelMessageSend("527763159000547330", "!add ")
//	}
//}

//func debounce(interval time.Duration, input chan int, f func(arg int)) {
//	var (
//		item int
//	)
//	for {
//		select {
//		case item = <-input:
//			fmt.Println("received a send on a spammy channel - might be doing a costly operation if not for debounce")
//		case <-time.After(interval):
//			f(item)
//		}
//	}
//}
