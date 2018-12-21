package f8

import (
	"regexp"

	"github.com/oct2pus/botutil/logging"
)

// ParseMod Parses modifier for fate die rolls
func ParseMod(i string) bool {
	compare, err := regexp.MatchString(
		"(\\+|-)?[0-9]*", i)
	logging.CheckError(err)

	if compare {
		return true
	}
	return false

}

// DieSymbol Dumb switch to convert -1 to 1 to a fate die symbol
func DieSymbol(i int64) string {
	switch i {
	case int64(-1):
		return "-"
	case int64(0):
		return "0"
	case int64(1):
		return "+"
	default:
		return "Oh gog." // something went wrong here
	}
}
