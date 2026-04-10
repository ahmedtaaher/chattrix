package utils

import "regexp"

var mentionRegex = regexp.MustCompile(`@(\w+)`)

func ExtractMentions(text string) []string {
	matches := mentionRegex.FindAllStringSubmatch(text, -1)

	var usernames []string
	for _, m := range matches {
		if len(m) > 1 {
			usernames = append(usernames, m[1])
		}
	}

	return usernames
}