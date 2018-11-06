package roll

import (
	"github.com/oct2pus/vriskabot/util/logging"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

type roll struct {
	Amount int64 // Amount of die being rolled
	Size   int64 // Size of die being rolled
	Mod    int64 // Modifier to roll
}

func New(a int64, b int64, c int64) roll {
	return roll{a, b, c}
}

// returns a series of random numbers (determined by die.sizeOfDie) in an int64
// slice, which is as large as die.numberOfDie
func RollTable(rolls roll) []int64 {
	var table []int64
	seed := time.Now()

	r := rand.New(rand.NewSource(seed.Unix()))

	for int64(len(table)) < rolls.Amount {
		table = append(table, (r.Int63n(rolls.Size) + 1))
	}

	return table

}

// breaks the roll into a string slice
// code assumes you've checked input prior
// TODO: add a break/error state
func DiceSlice(input string) []string {
	// [0] is the number of dice being rolled
	// [1] is the type of die
	// [2] is the modifier direction (positive/negative)
	// [3] is the size of the modifier (0 if none)

	divider := regexp.MustCompile("[0-9]+|[\\+|-]")

	dieSlice := divider.FindAllString(input, -1)

	if len(dieSlice) <= 2 {
		dieSlice = append(dieSlice, "+")
	}
	if len(dieSlice) <= 3 {
		dieSlice = append(dieSlice, "0")
	}
	return dieSlice
}

// turns the dieSlice string slice into a dieRoll object
func FromStrings(diceSlice []string) roll {
	var die roll
	var err error

	die.Amount, err = strconv.ParseInt(diceSlice[0], 0, 0)
	logging.CheckError(err)
	die.Size, err = strconv.ParseInt(diceSlice[1], 0, 0)
	logging.CheckError(err)
	die.Mod, err = strconv.ParseInt(diceSlice[3], 0, 0)
	logging.CheckError(err)

	// if number is negative is negative
	if diceSlice[2] == "-" {
		die.Mod = 0 - die.Mod
	}

	return die

}
