//Check Line 259 Dummy

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/oct2pus/vriskabot/etc"
	"github.com/oct2pus/vriskabot/logging"
	"github.com/oct2pus/vriskabot/parse"
	"github.com/oct2pus/vriskabot/roll"
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

// 'global' variables
var (
	// command line argument
	Token string

	self *discordgo.User
)

// initalize variables
func init() {
	executable, e := os.Executable()
	if e != nil {
		panic(e)
	}
	// command line argument
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

	// error logging
	logging.CreateLog()
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

	fmt.Println("Guilds: ", len(discordReady.Guilds))
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
	if message[0] == prefix && len(message) > 1 {
		switch message[1] {
		case "roll", "lroll", "hroll":
			if len(message) > 2 {
				embed, err := sendRoll(message[2], message[1])
				if !logging.CheckError(err) {
					discordSession.ChannelMessageSend(discordMessage.ChannelID, "Rolling!!!!!!!!")
					discordSession.ChannelMessageSendEmbed(discordMessage.ChannelID, embed)
				} else {
					discordSession.ChannelMessageSend(discordMessage.ChannelID,
						err.Error())
				}
			} else {
				discordSession.ChannelMessageSend(discordMessage.ChannelID, "Roll what?")
			}
		case "fortune", "8ball":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"w8 a moment ::::)")
		case "stats":
			discordSession.ChannelMessageSend(discordMessage.ChannelID,
				"place(8e)holder?")
		case "fate", "f8":
			var embed *discordgo.MessageEmbed
			var err error

			if len(message) > 2 {
				embed, err = sendF8Roll(message[2])
			} else {
				embed, err = sendF8Roll("0")
			}
			if err != nil {
				discordSession.ChannelMessageSend(discordMessage.ChannelID,
					err.Error())
			} else {
				discordSession.ChannelMessageSend(discordMessage.ChannelID, "Rolling!!!!!!!!")
				discordSession.ChannelMessageSendEmbed(discordMessage.ChannelID, embed)
			}
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

// for sending dice rolls related to the game fate
func sendF8Roll(modifier string) (*discordgo.MessageEmbed, error) {

	f8 := roll.new(4, 3, 0)
	if modifier != "" && parse.CheckFormated(modifier, "(\\+|-)?[0-9]*") {
		mod, err := strconv.ParseInt(modifier, 10, 64)
		logging.CheckError(err)
		roll.Mod = mod
	} else if modifier != "" {
		return nil, errors.New("W8 what?")
	}

	table := roll.RollTable(f8)

	// fate rolls are actually -1 to 1, not 1 to 3
	for i, _ := range rolls {
		rolls[i] -= 2
	}

	var f8Rolls []string

	for _, ele := range rolls {
		f8Rolls = append(f8Rolls, toF8DieSymbol(ele))
	}

	total := strconv.FormatInt(getTotal(rolls), 10)

	dieImage := "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/dfate.png"

	return dieRollEmbed(f8Rolls, strconv.FormatInt(roll.modifier, 10), total, dieImage), nil
}

func parseF8Mod(i string) bool {
	compare, err := regexp.MatchString(
		"(\\+|-)?[0-9]*", i)
	logging.CheckError(err)

	if compare {
		return true
	}
	return false

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
				Value:  "Avatar By mjölk#8323 ( http://cosmic-rumpus.tumblr.com/ )\nEmojis by Dzuk#1671 ( https://noct.zone/ )",
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Disclaimer",
				Value:  "Vriska8ot uses **Mutant Standard Emoji** (https://mutant.tech)\n**Mutant Standard Emoji** are licensed under CC-BY-NC-SA 4.0 (https://creativecommons.org/licenses/by-nc-sa/4.0/)",
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
	if !parse.CheckFormatted(diceString,
		"[1-9]+[0-9]*d[1-9]+[0-9]*((\\+|-){1}[0-9]*)?") {
		valid = false
	}

	// This is called valid because the internet has made a fool of me.
	if valid {

		dieSlices := roll.DiceSlice(diceString)
		dice := roll.FromStrings(dieSlices)

		if dice.Amount > 20 {
			return nil, errors.New("Why would anyone ever need to roll that many dice?")
		}

		rollTable := roll.RollTable(die)
		var stringTable []string
		var result int64

		// get result
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

		result += dice.Mod

		//convert int slice to string slice

		for i, ele := range rollTable {
			stringTable[i] = strconv.FormatInt(ele, 10)
		}

		dieImage := determineDieImage(dice)

		embed := dieRollEmbed(stringTable,
			strconv.FormatInt(dice.Mod, 10), strconv.FormatInt(result, 10),
			dieImage)

		return embed, nil
	}

	return nil, errors.New("You gotta format it like this!\n`vriska: roll XdX(+/-X)`")
}

// I stopped at rewriting embeds so they can be multi use
func dieRollEmbed(rollTable []string, mod string, result string,
	dieImage string) *discordgo.MessageEmbed {

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
				Value:  mod,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Result",
				Value:  result,
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: dieImage,
		},
	}

	return embed
}

func toF8DieSymbol(i int64) string {
	switch i {
	case int64(-1):
		return "-"
	case int64(0):
		return "0"
	case int64(1):
		return "+"
	default:
		return "Oh gog."
	}
}

// Takes roll table and returns a
func formatRollTable(table []string) string {
	fieldValue := "`"
	for x := 0; x < len(table); x++ {
		if x%4 == 0 && x != 0 {
			fieldValue += "`\n`"
		}
		if x != 0 && x%4 != 0 {
			fieldValue += "\t"
		}
		fieldValue += "|" +
			toCenter(table[x]) + "|"
	}

	fieldValue += "`"

	return fieldValue
}

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

// determines if the diceString input is formatted properly
func isDiceMessageFormated(diceString string) bool {
	// todo: fix +- bullshit with regexp
	compare, err := regexp.MatchString(
		"[1-9]+[0-9]*d[1-9]+[0-9]*((\\+|-){1}[0-9]*)?", diceString)
	logging.CheckError(err)

	if compare {
		return true
	}
	return false
}
