package color

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
)

// PrefixError adds a red `Error:` prefix
var PrefixError = Red + "Error:" + Reset

// PrefixWarn adds a yellow `Warn:` prefix
var PrefixWarn = Yellow + "Warn:" + Reset

// PrefixInfo add a green `Info:` prefix
var PrefixInfo = Green + "Info:" + Reset

// BlueText renders given text in blue color
func BlueText(msg string) string {
	return Blue + msg + Reset
}
