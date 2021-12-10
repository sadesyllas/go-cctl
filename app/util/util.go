package util

import "regexp"

func UnquoteParsedStringValue(value string) string {
	return regexp.MustCompile(`^(?:<|\")|(?:>|\")$`).ReplaceAllString(value, "")
}
