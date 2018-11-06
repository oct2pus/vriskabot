package parse

import (
	"github.com/oct2pus/vriskabot/util/logging"
	"regexp"
)

func CheckFormatted(input string, rgxp string) bool {
	// todo: fix +- bullshit with regexp
	compare, err := regexp.MatchString(rgxp, input)
	logging.CheckError(err)

	if compare {
		return true
	}
	return false
}
