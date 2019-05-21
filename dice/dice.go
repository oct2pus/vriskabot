package dice

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Roll contains 3 int64 values
// Amount is how many dice are being rolled
// Size is how large a dice being rolled is
// Mod is the arthmatic modifier to the final tally
type Roll struct {
	Amount int64
	Size   int64
	Mod    int64
}

//TODO: further seperate packages, only things using a rollstruct should be here

// New generates a new Roll struct.
func New(a int64, b int64, c int64) Roll {
	return Roll{a, b, c}
}

// Table returns a series of random numbers (determined by die.sizeOfDie) in an
//  int64 slice the size of rolls.Amount.
func Table(rolls Roll) []int64 {
	var table []int64
	seed := time.Now()

	r := rand.New(rand.NewSource(seed.Unix()))

	for int64(len(table)) < rolls.Amount {
		table = append(table, (r.Int63n(rolls.Size) + 1))
	}

	return table

}

// Slice breaks the Roll into a string slice
// code assumes you've checked input prior
// TODO: add a break/error state.
func Slice(input string) []string {
	// [0] is the number of dice being rolled
	// [1] is the type of die
	// [2] is the modifier direction (positive/negative)
	// [3] is the size of the modifier (0 if none)

	divider := regexp.MustCompile("[0-9]+|[\\+|-]")

	slice := divider.FindAllString(input, -1)

	if len(slice) <= 2 {
		slice = append(slice, "+")
	}
	if len(slice) <= 3 {
		slice = append(slice, "0")
	}
	return slice
}

// FromStringSlice takes slice (a []string) and outputs a Roll struct.
func FromStringSlice(slice []string) (Roll, error) {
	var die Roll
	var err error

	die.Amount, err = strconv.ParseInt(slice[0], 0, 0)
	if err != nil {
		return Roll{}, err
	}
	die.Size, err = strconv.ParseInt(slice[1], 0, 0)
	if err != nil {
		return Roll{}, err
	}
	die.Mod, err = strconv.ParseInt(slice[3], 0, 0)
	if err != nil {
		return Roll{}, err
	}

	if slice[2] == "-" {
		die.Mod = 0 - die.Mod
	}

	return die, nil
}

// FormatTable takes and returns a 'table' formatted for an embed
func FormatTable(table []string) string {
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

// DieImage returns an URL as a string based on an int64 value
func DieImage(face int64) string {
	switch {
	case face <= 4:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/art/d4.png"
	case face <= 6:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/art/d6.png"
	case face <= 8:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/art/d8.png"
	case face <= 10:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/art/d10.png"
	case face <= 12:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/art/d12.png"
	default:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/art/d20.png"
	}
}

// GetTotal gets the total value from an int64 slice
func GetTotal(arr []int64) int64 {

	sum := int64(0)
	for x := 0; x < len(arr); x++ {
		sum += arr[x]
	}

	return sum
}

// GetHighest gets the highest value from an int64 slice
func GetHighest(arr []int64) int64 {
	highest := int64(0)
	for x := 0; x < len(arr); x++ {
		if highest < arr[x] {
			highest = arr[x]
		}
	}

	return highest
}

// GetLowest gets the lowest value from an int64 slice
func GetLowest(arr []int64) int64 {
	lowest := arr[0]
	for x := 1; x < len(arr); x++ {
		if lowest > arr[x] {
			lowest = arr[x]
		}
	}

	return lowest
}

// IsFormated determines if the diceString input is formatted properly
func IsFormated(diceString string) bool {
	// todo: fix +- bullshit with regexp
	compare, err := regexp.MatchString(
		"[1-9]+[0-9]*d[1-9]+[0-9]*((\\+|-){1}[0-9]*)?", diceString)
	if err != nil {
		return false
	}

	if compare {
		return true
	}
	return false
}

// RollEmbed provides a nice, clean view into the results
func RollEmbed(rollTable []string, mod string, result string,
	dieImage string) *discordgo.MessageEmbed {

	embed := &discordgo.MessageEmbed{
		Color: 0x005682,
		Type:  "Roooooooolling!",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Rolls",
				Value:  FormatTable(rollTable),
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

// toCenter centers text
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
