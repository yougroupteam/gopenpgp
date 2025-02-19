// Package internal contains internal methods and constants.
package internal

import (
	"strings"

	"github.com/yougroupteam/gopenpgp/v2/constants"
)

func CanonicalizeAndTrim(text string) string {
	lines := strings.Split(text, "\n")

	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t\r")
	}

	return strings.Join(lines, "\r\n")
}

// CreationTimeOffset stores the amount of seconds that a signature may be
// created in the future, to compensate for clock skew.
const CreationTimeOffset = int64(60 * 60 * 24 * 2)

// ArmorHeaders is a map of default armor headers.
var ArmorHeaders = map[string]string{
	"Version": constants.ArmorHeaderVersion,
	"Comment": constants.ArmorHeaderComment,
}
