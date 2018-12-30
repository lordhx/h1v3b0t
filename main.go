package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
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
	dg.AddHandler(handleAdd)

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
