package utils

import (
	"strconv"
	"strings"
)

func Emoji(s string) (string, error) {
    r, err := strconv.ParseInt(strings.TrimPrefix(s, "\\U"), 16, 32)
    return string(r), err
}
