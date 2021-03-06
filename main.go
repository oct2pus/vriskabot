package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oct2pus/vriskabot/command"

	"github.com/oct2pus/bocto"
)

const (
	pre = "8"
)

func main() {
	// initalize variables
	var token string
	var vriska bocto.Bot
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	err := vriska.New("Vriska8ot", "8", token,
		"Hiiiiiiii?\n8y the way my prefix is `"+pre+"`. "+
			"Not that you neeeeeeeeded to know or anything.", "::::?\nAre you 8egging me to `"+pre+"help` you?", 0x005682)
	if err != nil {
		fmt.Printf("%v can't login\nerror: %v\n", vriska.Name, err)
		return
	}
	// add commands and responses
	vriska = addCommands(vriska)
	// Event Handlers
	vriska.Session.AddHandler(vriska.ReadyEvent)
	vriska.Session.AddHandler(vriska.MessageCreate)
	vriska.DisablePhrases = true

	// Open Bot
	err = vriska.Session.Open()
	if err != nil {
		fmt.Printf("Error opening connection: %v\nDump bot info %v\n",
			err,
			vriska.String())
	}
	// wait for ctrl+c to close.
	signalClose := make(chan os.Signal, 1)

	signal.Notify(signalClose,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
		os.Kill)
	<-signalClose

	vriska.Session.Close()
}

func addCommands(bot bocto.Bot) bocto.Bot {
	// alphabetical order, shorter first
	bot.AddCommand("about", command.Credits)
	bot.AddCommand("commands", command.Help)
	bot.AddCommand("credits", command.Credits)
	bot.AddCommand("discord", command.Discord)
	bot.AddCommand("f8", command.F8)
	bot.AddCommand("fate", command.F8)
	bot.AddCommand("help", command.Help)
	bot.AddCommand("hroll", command.HRoll)
	bot.AddCommand("invite", command.Invite)
	bot.AddCommand("lroll", command.LRoll)
	bot.AddCommand("roll", command.Roll)

	return bot
}
