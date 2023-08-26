package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	outfile, _ = os.Create("rssdcbot.log")
	l          = log.New(outfile, "", log.LstdFlags|log.Lshortfile)
)

func main() {
	ConnectToDC()
}

var token = "YOUR-TOKEN"
var Dg *discordgo.Session

func initSession() {
	// Initilaze the Discord Session
	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("Couldn't start Discord session! Error: ", err)
		l.Fatalf("Couldn't start Discord session! Error: %s", err)

	}

	Dg = dg

}

func ConnectToDC() {
	initSession()
	go ParseMeetups()
	time.Sleep(time.Second * 1)

	//Register messageCreate func
	Dg.AddHandler(messageCreate)
	Dg.Identify.Intents = discordgo.IntentDirectMessageReactions

	//If the user press CTRL+C, exit.
	fmt.Println("Bot is running. To exit press CTRL+C")
	l.Println("[INFO] Bot is running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	Dg.Close()
	l.Println("[INFO] Bot stopped.")

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!rssbot" {
		s.ChannelMessageSend(m.ChannelID, "Hey! I'm here.")
		l.Printf("[INFO] %s called the bot.", m.Author)
	}
}
