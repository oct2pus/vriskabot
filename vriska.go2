package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/oct2pus/botutil/embed"
	"github.com/oct2pus/botutil/etc"
	"github.com/oct2pus/botutil/logging"
	"github.com/oct2pus/botutil/parse"
	"github.com/oct2pus/vriskabot/util/dice"
	"github.com/oct2pus/vriskabot/util/dice/f8"
)

// Prefix Const
const (
	prefix = "vriska:"
)

// contains all the information about a dice roll

// 'global' variables
var (
	Token string
	self  *discordgo.User
	color int
)

// initalize variables
func init() {
	// command line argument
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

	// error logging
	logging.CreateLog()

	color = 0x005682
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

	if logging.CheckError(err) {
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

// messageCreate will be called  every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(session *discordgo.Session,
	message *discordgo.MessageCreate) {

	messageSlice := etc.StringSlice(message.Message.Content)
	// Ignore all messages created by the bot itself
	if message.Author.Bot == true {
		return
	}

	// commands
	// TODO: clean this up; its so damn ugly
	// TODO: split how errors are reported internally and externally
	if messageSlice[0] == prefix && len(messageSlice) > 1 {
		switch messageSlice[1] {
		case "roll", "lroll", "hroll":
			if len(messageSlice) > 2 {
				embed, err := sendRoll(messageSlice[2], messageSlice[1])
				if !logging.CheckError(err) {
					session.ChannelMessageSend(message.ChannelID, "Rolling"+
						"!!!!!!!!")
					session.ChannelMessageSendEmbed(message.ChannelID, embed)
				} else {
					session.ChannelMessageSend(message.ChannelID,
						err.Error())
				}
			} else {
				session.ChannelMessageSend(message.ChannelID, "Roll what?")
			}
		case "fortune", "8ball":
			session.ChannelMessageSend(message.ChannelID,
				"w8 a moment ::::)")
		case "stats":
			session.ChannelMessageSend(message.ChannelID,
				"place(8e)holder?")
		case "fate", "f8":
			var embed *discordgo.MessageEmbed

			if len(messageSlice) > 2 {
				embed = sendF8Roll(messageSlice[2])
			} else {
				embed = sendF8Roll("0")
			}

			session.ChannelMessageSend(message.ChannelID, "Rolling!!!!!!!!")
			session.ChannelMessageSendEmbed(message.ChannelID, embed)
		case "discord":
			session.ChannelMessageSend(message.ChannelID,
				"https://discord.gg/PGVh2M8")
		case "invite":
			session.ChannelMessageSend(message.ChannelID,
				"<https://discordapp.com/oauth2/authorize?client_id=497943811"+
					"700424704&scope=bot&permissions=281600>")
		case "help", "commands":
			session.ChannelMessageSend(message.ChannelID,
				"My commands are:\n`roll`\n`lroll`\n`hroll`\n`f8`\n`discord`"+
					"\n`invite`\n`help`\n`about`")
		case "about", "credits":
			session.ChannelMessageSendEmbed(message.ChannelID,
				embed.CreditsEmbed("Vriska8ot", " mjölk#8323 "+
					"( http://cosmic-rumpus.tumblr.com/ )",
					"", "Emojis by Dzuk#1671 ( https://noct.zone/ )",
					0x005682))
		default:
			session.ChannelMessageSend(message.ChannelID,
				"::::?")
		}
	}

	for _, ele := range message.Mentions {
		if ele.Username == self.Username {
			session.ChannelMessageSend(message.ChannelID,
				"Hiiiiiiii?\n8y the way my prefix is '`vriska: `'. "+
					"Not that you neeeeeeeeded to know or anything.")
		}
	}
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

		dieSlices := dice.Slice(diceString)
		die := dice.FromStringSlice(dieSlices)

		if die.Amount > 20 {
			return nil, errors.New("Why would anyone ever need to roll that " +
				"many dice?")
		}

		rollTable := dice.Table(die)
		var stringTable []string
		var result int64

		// get result
		switch commandInput {
		case "roll":
			result = dice.GetTotal(rollTable)
		case "lroll":
			result = dice.GetLowest(rollTable)
		case "hroll":
			result = dice.GetHighest(rollTable)
		case "default":
			return nil, errors.New("Holy sh8t dont break me!")
		}

		result += die.Mod

		//convert int slice to string slice

		for _, ele := range rollTable {
			stringTable = append(stringTable, strconv.FormatInt(ele, 10))
		}

		dieImage := dice.DieImage(die.Size)

		embed := dice.RollEmbed(stringTable,
			strconv.FormatInt(die.Mod, 10), strconv.FormatInt(result, 10),
			dieImage)

		return embed, nil
	}

	return nil, errors.New("You gotta format it like this!\n`vriska: " +
		"roll XdX(+/-X)`")
}

// sendF8Roll sends dice rolls related to the game "Fate"
func sendF8Roll(modifier string) *discordgo.MessageEmbed {

	var mod int64
	var err error
	if modifier != "" && parse.CheckFormatted(modifier, "(\\+|-)?[0-9]+$") {
		mod, err = strconv.ParseInt(modifier, 10, 64)
		logging.CheckError(err)
	} else {
		mod = 0
		err = nil
	}
	die := dice.New(4, 3, mod)

	table := dice.Table(die)

	// fate rolls are actually -1 to 1, not 1 to 3
	for i := range table {
		table[i] -= 2
	}

	var rolls []string

	for _, ele := range table {
		rolls = append(rolls, f8.DieSymbol(ele))
	}

	total := strconv.FormatInt(dice.GetTotal(table)+die.Mod, 10)

	dieImage := "https://raw.githubusercontent.com/oct2pus/vriskabot/master/e" +
		"moji/dfate.png"

	return dice.RollEmbed(rolls, strconv.FormatInt(die.Mod, 10), total,
		dieImage)
}
