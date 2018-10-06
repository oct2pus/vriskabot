package main

import (
	"errors"
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

// contains all the information about a dice roll
type dieRoll struct {
	numberOfDie int64
	sizeOfDie   int64
	modifier    int64
}

// 'global' variables
var (
	// command line argument
	Token string
	// error logging
	Log         *log.Logger
	currentTime string
	self        *discordgo.User
)

// initalize variables
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
	file, err := os.Create(path + ".logs@" + currentTime + ".log")
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.Ldate|log.Ltime|log.Llongfile|log.LUTC)
}

// buddy its main
func main() {

	// Create a new Discord session using the provided bot token.
	// token must be prefaced with "Bot "
	bot, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Bot Event Handlers
	bot.AddHandler(messageCreate)
	bot.AddHandler(ready)

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

// This function is called when the bot connects to discord
func ready(discordSession *discordgo.Session, discordReady *discordgo.Ready) {
	discordSession.UpdateStatus(0, "prefix: \""+prefix+" \"")
	self = discordReady.User
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
			if len(message) > 2 {
				embed, err := sendRoll(message[2], message[1])
				if !checkError(err) {
					discordSession.ChannelMessageSend(discordMessage.ChannelID, "Rolling!!!!!!!!")
					discordSession.ChannelMessageSendEmbed(discordMessage.ChannelID, embed)
				} else {
					discordSession.ChannelMessageSend(discordMessage.ChannelID,
						err.Error())
				}
			} else {
				discordSession.ChannelMessageSend(discordMessage.ChannelID, "Roll what?")
			}
		case "stats":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"placeholder?")
		case "fate":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"placeholder?!")
		case "discord":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"https://discord.gg/PGVh2M8")
		case "invite":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"<https://discordapp.com/oauth2/authorize?client_id=497943811700424704&scope=bot&permissions=281600>")
		case "help", "commands":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"My help is currently incomplete, 8ut my commands are:\n`roll`\n`lroll`\n`hroll`\n`discord`\n`invite`\n`help`\n`about`")
		case "about", "credits":
			discordSession.ChannelMessageSendEmbed(discordMessage.ChannelID, getCredits())
		default:
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"::::?")
		}
	}

	for _, ele := range discordMessage.Mentions {
		if ele.Username == self.Username {
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"Hiiiiiiii?\n8y the way my prefix is '`vriska: `'. Not that you neeeeeeeeded to know or anything.")
		}
	}
}

func getCredits() *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color: 0x005682,
		Type:  "A8out",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Vriska8ot",
				Value:  "Created by \\🐙\\🐙#0413 ( http://oct2pus.tumblr.com/ )\nVriska8ot uses the 'discordgo' library\n( https://github.com/bwmarrin/discordgo/ )",
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Special Thanks",
				Value:  "Avatar By mjölk#8323 ( http://cosmic-rumpus.tumblr.com/ )\nEmojis by Dzuk#1671",
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Disclaimer",
				Value:  "Vriska8ot uses **Mutant Standard Emoji** (https://mutant.tech)\n**Mutant Standard Emoji** are licensed under CC-BY-NC-SA 4.0 (https://creativecommons.org/licenses/by-nc-sa/4.0/) ",
				Inline: false,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://raw.githubusercontent.com/oct2pus/vriskabot/master/art/avatar.png",
		},
	}

	return embed
}

// performs the 'math' for a roll, lroll, or hroll function, returns a
// MessageEmbed
func sendRoll(diceString string, commandInput string) (*discordgo.MessageEmbed,
	error) {

	valid := true
	if !isDiceMessageFormated(diceString) {
		valid = false
	}

	if valid {

		dieSlices := divideIntoDieSlices(diceString)
		die := convertToDieRollStruct(dieSlices)
		rollTable := determineRollTable(die)

		if die.numberOfDie > 20 {
			return nil, errors.New("Why would anyone ever need to roll that many dice?")
		}

		var result int64

		switch commandInput {
		case "roll":
			result = getTotal(rollTable)
		case "lroll":
			result = getLowest(rollTable)
		case "hroll":
			result = getHighest(rollTable)
		case "default":
			return nil, errors.New("Holy sh8t don't 8reak me!!!!!!!!")
		}

		result += die.modifier

		dieImage := determineDieImage(die)

		embed := &discordgo.MessageEmbed{
			Color: 0x005682,
			Type:  "Roooooooolling!",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Rolls",
					Value:  formatRollTable(rollTable),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Modifier",
					Value:  strconv.FormatInt(die.modifier, 10),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Result",
					Value:  strconv.FormatInt(result, 10),
					Inline: true,
				},
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: dieImage,
			},
		}

		return embed, nil
	}
	return nil, errors.New("You gotta format it like this!\n`vriska: roll XdX(+/-X)`")
}

func formatRollTable(table []int64) string {
	fieldValue := "`"
	for x := 0; x < len(table); x++ {
		if x%4 == 0 && x != 0 {
			fieldValue += "`\n`"
		}
		if x != 0 && x%4 != 0 {
			fieldValue += "\t"
		}
		fieldValue += "|" +
			toCenter(strconv.FormatInt(table[x], 10)) + "|"
	}

	fieldValue += "`"

	return fieldValue
}

// centers text
// im doing this the shitty not expandable way because ive been defeated
func toCenter(s string) string {
	switch len(s) {
	case 1:
		return " " + s + " "
	case 2:
		return " " + s
	default:
		return s
	}
}

/*
// centers text, properly, but for some reason throws a hissy fit if i use
// spaces
func toCenter(s string, i int) string {
	if i > len(s) {
		o := i - len(s)
		ns := spaceLoop("", o) + s

		return ns
	}
	return s
}

// adds 'i' spaces to string 's'
func spaceLoop(s string, i int) string {

	for len(s) < i {
		s += "_"
	}
	return s
}
*/
// determines what image to use
func determineDieImage(die dieRoll) string {
	switch {
	case die.sizeOfDie <= 4:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d4.png"
	case die.sizeOfDie <= 6:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d6.png"
	case die.sizeOfDie <= 8:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d8.png"
	case die.sizeOfDie <= 10:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d10.png"
	case die.sizeOfDie <= 12:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d12.png"
	case die.sizeOfDie > 12:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d20.png"
	}
	return ""
}

// gets the total value from an int64 slice
func getTotal(arr []int64) int64 {

	sum := int64(0)
	for x := 0; x < len(arr); x++ {
		sum += arr[x]
	}

	return sum
}

// gets the highest value from an int64 slice
func getHighest(arr []int64) int64 {
	highest := int64(0)
	for x := 0; x < len(arr); x++ {
		if highest < arr[x] {
			highest = arr[x]
		}
	}

	return highest
}

// gets the lowest value from an int64 slice
func getLowest(arr []int64) int64 {
	lowest := arr[0]
	for x := 1; x < len(arr); x++ {
		if lowest > arr[x] {
			lowest = arr[x]
		}
	}

	return lowest
}

// returns a series of random numbers (determined by die.sizeOfDie) in an int64
// slice, which is as large as die.numberOfDie
func determineRollTable(die dieRoll) []int64 {
	var rolls []int64
	seed := time.Now()

	r := rand.New(rand.NewSource(seed.Unix()))

	for int64(len(rolls)) < die.numberOfDie {
		rolls = append(rolls, (r.Int63n(die.sizeOfDie) + 1))
	}

	return rolls

}

// determines if the diceString input is formatted properly
func isDiceMessageFormated(diceString string) bool {
	// todo: fix +- bullshit with regexp
	compare, err := regexp.MatchString(
		"[1-9]+[0-9]*d[1-9]+[0-9]*((\\+|-){1}[0-9]*)?", diceString)
	checkError(err)

	if compare {
		return true
	}
	return false
}

// breaks the dieString into a string slice
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

// turns the dieSlice string slice into a dieRoll object
func convertToDieRollStruct(dieSlice []string) dieRoll {
	var die dieRoll
	var err error
	die.numberOfDie, err = strconv.ParseInt(dieSlice[0], 0, 0)
	checkError(err)
	die.sizeOfDie, err = strconv.ParseInt(dieSlice[1], 0, 0)
	checkError(err)

	die.modifier, err = strconv.ParseInt(dieSlice[3], 0, 0)
	checkError(err)

	// if number is negative is negative
	if dieSlice[2] == "-" {
		die.modifier = 0 - die.modifier
	}

	return die

}

// logs errors
func checkError(err error) bool {
	if err != nil {
		fmt.Println("error: ", err)
		Log.Println("error: ", err)
		return true
	}
	return false
}

// converts text to lowercase substrings
func parseText(m string) []string {

	m = strings.ToLower(m)
	return strings.Split(m, " ")
}
