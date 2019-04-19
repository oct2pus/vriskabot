package command

import (
	"regexp"
	"strconv"

	"github.com/oct2pus/bot/bot"

	"vriskabot/dice/f8"

	"github.com/oct2pus/bot/embed"

	"vriskabot/dice"

	"github.com/bwmarrin/discordgo"
)

// F8 represents a F8 dice rice.
func F8(bot bot.Bot, message *discordgo.MessageCreate, input []string) {
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

	go embed.SendEmbededMessage(bot.Session, message.ChannelID,
		dice.RollEmbed(rolls, strconv.FormatInt(die.Mod, 10), total,
			dieImage))
}

// Discord posts my discord URL.
func Discord(bot bot.Bot, message *discordgo.MessageCreate, input []string) {
	go embed.SendMessage(bot.Session, message.ChannelID,
		"https://discord.gg/PGVh2M8")
}

// Invite posts a link to invite Vriska8ot.
func Invite(bot bot.Bot, message *discordgo.MessageCreate, input []string) {
	go embed.SendMessage(bot.Session, message.ChannelID,
		"<https://discordapp.com/oauth2/authorize?client_id=497943811"+
			"700424704&scope=bot&permissions=281600>")
}

// Credits accreditates users for their contributions.
func Credits(bot bot.Bot,
	message *discordgo.MessageCreate,
	input []string) {
	go embed.SendEmbededMessage(bot.Session, message.ChannelID,
		embed.CreditsEmbed(bot.Name,
			"(milk wizard#8323 http://cosmic-rumpus.tumblr.com/ )",
			"",
			"Dzuk#1671 ( https://noct.zone/ )",
			"https://raw.githubusercontent.com/oct2pus/vriskabot/master/art/"+
				"vriskabot.png",
			bot.Color))
}

// Help returns a list of commands.
func Help(bot bot.Bot,
	message *discordgo.MessageCreate,
	input []string) {

	go embed.SendMessage(bot.Session, message.ChannelID,
		"My commands are:\n`roll`\n`lroll`\n`hroll`\n`f8`\n`discord`"+
			"\n`invite`\n`help`\n`about`")
}

// Roll returns a normal type of roll.
func Roll(bot bot.Bot, message *discordgo.MessageCreate, input []string) {
	if noInput(input) {
		input = append(input, "")
	}
	go roll(bot, message, input[0], "roll")
}

// LRoll returns the lowest die in a roll.
func LRoll(bot bot.Bot, message *discordgo.MessageCreate, input []string) {
	if noInput(input) {
		input = append(input, "")
	}
	go roll(bot, message, input[0], "lroll")
}

// HRoll returns the highest die in a roll.
func HRoll(bot bot.Bot, message *discordgo.MessageCreate, input []string) {
	if noInput(input) {
		input = append(input, "")
	}
	go roll(bot, message, input[0], "hroll")
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
func roll(bot bot.Bot,
	message *discordgo.MessageCreate,
	diceString, com string) {

	if !checkFormatted(diceString,
		"[1-9]+[0-9]*d[1-9]+[0-9]*((\\+|-){1}[0-9]*)?") {
		go embed.SendMessage(bot.Session, message.ChannelID,
			"You gotta format it like this!\n`vriska: "+
				"roll XdX(+/-X)`")
		return
	}

	dieSlices := dice.Slice(diceString)
	die, err := dice.FromStringSlice(dieSlices)
	if err != nil {
		go embed.SendMessage(bot.Session, message.ChannelID, "That num8er is"+
			" waaaaaaaay to 8ig for me the handle.")
		return
	}

	if die.Amount > 20 {
		go embed.SendMessage(bot.Session, message.ChannelID,
			"Why would anyone ever need to roll that "+
				"many dice?")
		return
	}

	embed.SendMessage(bot.Session, message.ChannelID,
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
		go embed.SendMessage(bot.Session, message.ChannelID,
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

	go embed.SendEmbededMessage(bot.Session, message.ChannelID, emb)
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
