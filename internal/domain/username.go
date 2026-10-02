package domain

import "strings"

func NormalizeTelegramUsername(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "@")
	return strings.ToLower(s)
}
