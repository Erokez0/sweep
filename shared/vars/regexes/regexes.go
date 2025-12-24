package regexes

import "regexp"

var (
	ColorRegex = regexp.MustCompile(`^(#([A-Fa-f0-9]{2,6}|[A-Fa-f0-9]{6})|(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9]))$`)
	KeyRegex = regexp.MustCompile(`^((ctrl\+)?(tab\+)?(([a-z])|enter))$`)
)
