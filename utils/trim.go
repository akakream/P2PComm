package utils

import "unicode/utf8"

func TrimTheSlashInTheBeginning(key string) string {
        c, i := utf8.DecodeRuneInString(key)
        if c == []rune("/")[0] {
            return key[i:]
        }
        return key
}
