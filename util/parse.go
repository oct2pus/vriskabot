package parse

import (
	"./logging"
	"regexp"
)

func CheckFormated(input string, regexp string) bool {
	// todo: fix +- bullshit with regexp
	compare, err := regexp.MatchString(regexp, diceString)
	logging.CheckError(err)

	if compare {
		return true
	}
	return false
}
