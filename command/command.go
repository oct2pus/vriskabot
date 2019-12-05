package command

import (
	"regexp"
	"strconv"

	"github.com/oct2pus/bocto"

	"github.com/oct2pus/vriskabot/dice/f8"

	"github.com/oct2pus/vriskabot/dice"

	"github.com/bwmarrin/discordgo"
)

// Credits provides attributation.
func Credits(bot bocto.Bot,
	message *discordgo.MessageCreate,
	input []string) {
	
	bot.Session.ChannelMessageSendEmbed(message.ChannelID,
		bocto.CreditsEmbed(
			bot.Name,
			bot.Self.AvatarURL(""),
			bot.Color,
			true,
			bocto.Contributor{
				Name:		"\\🐙\\🐙",
				URL:		"https://oct2pus.tumblr.com/",
				Message: 	"**Developed** by **%v** (%v)",
				Type: 		"Developer",
			},
			bocto.Contributor{
				Name:		"Discordgo",
				URL:		"https://github.com/bwmarrin/discordgo/",
				Message:	"**Vriska8ot** uses the **%v** library (%v)",
				Type:		"Library",
			},
			bocto.Contributor{
				Name:		"Dzuk",
				URL: 		"https://noct.zone/",
				Message: 	"**Emoji** by **%v** (%v)",
				Type:		"Artist",
			},
			bocto.Contributor{
				Name:		"Milk Wizard",
				URL: 		"https://mi1k-wizard.tumblr.com/",
				Message: 	"**Avatar** by **%v** (%v)",
				Type: 		"Artist",
			},
		),
	)
}

// Discord posts my discord URL.
func Discord(bot bocto.Bot, message *discordgo.MessageCreate, input []string) {
	bot.Session.ChannelMessageSend(message.ChannelID,
		"https://discord.gg/PFCGhJQ")
}

// F8 represents a F8 dice rice.
func F8(bot bocto.Bot, message *discordgo.MessageCreate, input []string) {
	var mod int64
	var err error
	if len(input) != 0 && checkFormatted(input[0], "(\\+|-)?[0-9]+$") {
		mod, err = strconv.ParseInt(input[0], 10, 64)
		if err != nil {
			return
		}
	} else {
		mod = 0
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

	dieImage := "https://raw.githubusercontent.com/oct2pus/vriskabot/master/" +
		"art/dfate.png"

	bot.Session.ChannelMessageSendEmbed(message.ChannelID,
		dice.RollEmbed(rolls, strconv.FormatInt(die.Mod, 10), total,
			dieImage))
}

// Help returns a list of commands.
func Help(bot bocto.Bot,
	message *discordgo.MessageCreate,
	input []string) {

	bot.Session.ChannelMessageSend(message.ChannelID,
		"My commands are:\n`roll`\n`lroll`\n`hroll`\n`f8`\n`discord`"+
			"\n`invite`\n`help`\n`about`")
}

// HRoll returns the highest die in a roll.
func HRoll(bot bocto.Bot, message *discordgo.MessageCreate, input []string) {
	if noInput(input) {
		input = append(input, "")
	}
	go roll(bot, message, input[0], "hroll")
}

// Invite posts a link to invite Vriska8ot.
func Invite(bot bocto.Bot, message *discordgo.MessageCreate, input []string) {
	bot.Session.ChannelMessageSend(message.ChannelID,
		"<https://discordapp.com/oauth2/authorize?client_id=497943811"+
			"700424704&scope=bot&permissions=281600>")
}

// LRoll returns the lowest die in a roll.
func LRoll(bot bocto.Bot, message *discordgo.MessageCreate, input []string) {
	if noInput(input) {
		input = append(input, "")
	}
	go roll(bot, message, input[0], "lroll")
}

// Roll returns a normal type of roll.
func Roll(bot bocto.Bot, message *discordgo.MessageCreate, input []string) {
	if noInput(input) {
		input = append(input, "")
	}
	go roll(bot, message, input[0], "roll")
}

func checkFormatted(input string, rgxp string) bool {
	// todo: fix +- bullshit with regexp
	compare, err := regexp.MatchString(rgxp, input)
	if err != nil {
		return false
	}

	if compare {
		return true
	}
	return false
}

// noInput catches a gotcha, i[] can be len = 0, return true if 0 or smaller.
func noInput(i []string) bool {
	if len(i) <= 0 {
		return true
	}
	return false
}

// roll performs the 'math' for a roll, lroll, or hroll function, should not
// be accessed directly.
func roll(bot bocto.Bot,
	message *discordgo.MessageCreate,
	diceString, com string) {

	if !checkFormatted(diceString,
		"[1-9]+[0-9]*d[1-9]+[0-9]*((\\+|-){1}[0-9]*)?") {
		bot.Session.ChannelMessageSend(message.ChannelID,
			"You gotta format it like this!\n`vriska: "+
				"roll XdX(+/-X)`")
		return
	}

	dieSlices := dice.Slice(diceString)
	die, err := dice.FromStringSlice(dieSlices)
	if err != nil {
		bot.Session.ChannelMessageSend(message.ChannelID, "That num8er is"+
			" waaaaaaaay to 8ig for me the handle.")
		return
	}

	if die.Amount > 20 {
		bot.Session.ChannelMessageSend(message.ChannelID,
			"Why would anyone ever need to roll that "+
				"many dice?")
		return
	}

	bot.Session.ChannelMessageSend(message.ChannelID,
		"Rolling!!!!!!!!")

	rollTable := dice.Table(die)
	var stringTable []string
	var result int64

	switch com {
	case "roll":
		result = dice.GetTotal(rollTable)
	case "lroll":
		result = dice.GetLowest(rollTable)
	case "hroll":
		result = dice.GetHighest(rollTable)
	case "default": // something REALLY bad happened if this is reached
		bot.Session.ChannelMessageSend(message.ChannelID,
			"Holy sh8t dont break me!!!!!!!!")
		return
	}
	result += die.Mod

	//convert int slice to string slice

	for _, ele := range rollTable {
		stringTable = append(stringTable, strconv.FormatInt(ele, 10))
	}

	dieImage := dice.DieImage(die.Size)

	emb := dice.RollEmbed(stringTable,
		strconv.FormatInt(die.Mod, 10), strconv.FormatInt(result, 10),
		dieImage)

	bot.Session.ChannelMessageSendEmbed(message.ChannelID, emb)
}
