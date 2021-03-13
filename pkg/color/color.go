package color

var (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	blue   = "\033[1;34m"
)

// PrefixError adds a red `Error:` prefix
var PrefixError = red + "ERROR:" + reset

// PrefixWarn adds a yellow `Warn:` prefix
var PrefixWarn = yellow + "WARN:" + reset

// PrefixInfo add a green `Info:` prefix
var PrefixInfo = green + "INFO:" + reset

// BlueText renders given text in blue color
func BlueText(msg string) string {
	return blue + msg + reset
}
