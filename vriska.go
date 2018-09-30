package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Prefix Const
const (
	prefix = "vriska:"
)

var (
	// command line argument
	Token string
	// error logging
	Log         *log.Logger
	currentTime string
)

func init() {
	executable, e := os.Executable()
	if e != nil {
		panic(e)
	}
	path := filepath.Dir(executable)

	// command line argument
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	// error logging
	currentTime = time.Now().Format("2006-01-02@15h04m")
	file, err := os.Create(path + "logs@" + currentTime + ".log")
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.Ldate|log.Ltime|log.Llongfile|log.LUTC)
}

func main() {

	// Create a new Discord session using the provided bot token.
	bot, err := discordgo.New("Bot " + Token) // token must be prefaced with "Bot "
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	bot.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		Log.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(discordSession *discordgo.Session,
	discordMessage *discordgo.MessageCreate) {

	message := parseText(discordMessage.Message.Content)
	// Ignore all messages created by the bot itself
	if discordMessage.Author.Bot == true {
		return
	}
	// commands
	if message[0] == prefix {
		switch message[1] {
		case "roll":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				strconv.FormatBool(isDiceMessageFormated(message[2])))
			fmt.Println(convertToDiceArray(message[2]))
		case "lroll":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"placeholder!")
		case "hroll":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"placeholder!")
		case "stats":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"placeholder?")
		case "fate":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"placeholder?!")
		case "discord":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"https://discord.gg/PGVh2M8")
		default:
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"::::?")
		}
	}

	// text responses
	// ... but nobody came
}

func isDiceMessageFormated(diceString string) bool {
	compare, err := regexp.MatchString("[0-9]+d[0-9]+((\\+|-)[0-9])?", diceString)
	if err != nil {
		Log.Println("error, ", err)
	}
	if compare {
		return true
	} else {
		return false
	}
}

func convertToDiceArray(diceString string) []int {
	// input[0] is the number of dice being rolled
	// input[1] is the type of die
	// input[2] is the size of the modifier (0 if none)
	// input[3] and above are irrelevant

	divider := regexp.MustCompile("(d|\\+|-)")
	parsedDiceString := divider.FindStringSubmatch(diceString)

	var DiceStringAsIntArray []int

	for i, ele := range parsedDiceString {
		x, err := strconv.ParseInt(ele, 0, 0)
		DiceStringAsIntArray[i] = int(x)
		if err != nil {
			fmt.Println(err)
			Log.Println(err)
		}
	}

	return DiceStringAsIntArray

}

// converts text to lowercase substrings
func parseText(m string) []string {

	m = strings.ToLower(m)
	return strings.Split(m, " ")
}
