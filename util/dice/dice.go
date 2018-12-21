package dice

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/oct2pus/botutil/etc"
	"github.com/oct2pus/botutil/logging"
)

type roll struct {
	Amount int64 // Amount of die being rolled
	Size   int64 // Size of die being rolled
	Mod    int64 // Modifier to roll
}

// New generates a new roll struct.
func New(a int64, b int64, c int64) roll {
	return roll{a, b, c}
}

// Table returns a series of random numbers (determined by die.sizeOfDie) in an int64
// slice the size of rolls.Amount.
func Table(rolls roll) []int64 {
	var table []int64
	seed := time.Now()

	r := rand.New(rand.NewSource(seed.Unix()))

	for int64(len(table)) < rolls.Amount {
		table = append(table, (r.Int63n(rolls.Size) + 1))
	}

	return table

}

// Slice breaks the roll into a string slice
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

// FromStringSlice takes slice (a []string) and outputs a roll struct.
func FromStringSlice(slice []string) roll {
	var die roll
	var err error

	die.Amount, err = strconv.ParseInt(slice[0], 0, 0)
	logging.CheckError(err)
	die.Size, err = strconv.ParseInt(slice[1], 0, 0)
	logging.CheckError(err)
	die.Mod, err = strconv.ParseInt(slice[3], 0, 0)
	logging.CheckError(err)

	// if number is negative is negative
	if slice[2] == "-" {
		die.Mod = 0 - die.Mod
	}

	return die

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
			etc.ToCenter(table[x]) + "|"
	}

	fieldValue += "`"

	return fieldValue
}

// DieImage returns an URL as a string based on an int64 value
func DieImage(face int64) string {
	switch {
	case face <= 4:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d4.png"
	case face <= 6:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d6.png"
	case face <= 8:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d8.png"
	case face <= 10:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d10.png"
	case face <= 12:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d12.png"
	default:
		return "https://raw.githubusercontent.com/oct2pus/vriskabot/master/emoji/d20.png"
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
