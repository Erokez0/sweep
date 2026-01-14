package regexes

import "regexp"

const (
	colorPattern = `^(#([A-Fa-f0-9]{2,6}|[A-Fa-f0-9]{6})|(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9]))$`
	keyPressPattern = `^((6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5][0-9]{4}|[0-9]{1,4})?(alt\+)?(ctrl\+)?(shift\+)?(\w|\\|enter|pgup|pgdn|delete|backspace|tab|~|\$|#|@)+)$`
)

var (
	ColorRegex = regexp.MustCompile(colorPattern)
	KeyPressRegex = regexp.MustCompile(keyPressPattern)
)
