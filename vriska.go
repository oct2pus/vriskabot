package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
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

type dieRoll struct {
	numberOfDie  int64
	sizeOfDie    int64
	modDirection bool // true positive | false negative
	modifier     int64
}

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
		case "roll", "lroll", "hroll":
			returnRoll(message[2], discordSession, discordMessage.ChannelID, message[1])
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

func returnRoll(diceString string, discordSession *discordgo.Session, channelID string, commandInput string) {
	valid := true
	if !isDiceMessageFormated(diceString) {
		valid = false
	}

	if valid {

		dieSlices := divideIntoDieSlices(diceString)
		die := convertToDieRollStruct(dieSlices)
		rollTable := determineRollTable(die)

		var result int64

		switch commandInput {
		case "roll":
			result = getTotal(rollTable)
		case "lroll":
			result = getLowest(rollTable)
		case "hroll":
			result = getHighest(rollTable)
		case "default":
			discordSession.ChannelMessageSend(channelID, "Something went hoooooooorribly wrong!!!!!!!! ::::(")
		}

	} else {
		discordSession.ChannelMessageSend(channelID, "::::?")
	}
}

func getTotal(arr []int64) int64 {

	sum := int64(0)
	for x := 0; x < len(arr); x++ {
		sum += arr[x]
	}

	return sum
}

func getHighest(arr []int64) int64 {
	highest := int64(0)
	for x := 0; x < len(arr); x++ {
		if highest < arr[x] {
			highest = arr[x]
		}
	}

	return highest
}

func getLowest(arr []int64) int64 {
	lowest := arr[0]
	for x := 1; x < len(arr); x++ {
		if lowest > arr[x] {
			lowest = arr[x]
		}
	}

	return lowest
}

func determineRollTable(die dieRoll) []int64 {
	var rolls []int64
	seed := time.Now()

	r := rand.New(rand.NewSource(seed.Unix()))

	for int64(len(rolls)) < die.numberOfDie {
		rolls = append(rolls, (r.Int63n(die.sizeOfDie) + 1))
	}

	return rolls

}

func isDiceMessageFormated(diceString string) bool {
	compare, err := regexp.MatchString("[1-9]+[0-9]*d[1-9]+[0-9]*((\\+|-){1}[0-9]*)?", diceString) // todo: fix +- bullshit

	checkError(err)

	if compare {
		return true
	}
	return false
}

func divideIntoDieSlices(dieString string) []string {
	// [0] is the number of dice being rolled
	// [1] is the type of die
	// [2] is the modifier direction (positive/negative)
	// [3] is the size of the modifier (0 if none)

	divider := regexp.MustCompile("[0-9]+|[\\+|-]")

	dieSlice := divider.FindAllString(dieString, -1)

	if len(dieSlice) <= 2 {
		dieSlice = append(dieSlice, "+")
	}
	if len(dieSlice) <= 3 {
		dieSlice = append(dieSlice, "0")
	}
	return dieSlice
}

func convertToDieRollStruct(dieSlice []string) dieRoll {
	var die dieRoll
	var err error
	die.numberOfDie, err = strconv.ParseInt(dieSlice[0], 0, 0)
	checkError(err)
	die.sizeOfDie, err = strconv.ParseInt(dieSlice[1], 0, 0)
	checkError(err)
	if dieSlice[2] != "-" {
		die.modDirection = true
	} else {
		die.modDirection = false
	}
	checkError(err)
	die.modifier, err = strconv.ParseInt(dieSlice[3], 0, 0)
	checkError(err)

	return die

}

func checkError(err error) {
	if err != nil {
		fmt.Println("error: ", err)
		Log.Println("error: ", err)
	}
}

// converts text to lowercase substrings
func parseText(m string) []string {

	m = strings.ToLower(m)
	return strings.Split(m, " ")
}
