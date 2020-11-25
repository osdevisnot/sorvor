package logger

import (
	"log"
	"strings"

	"github.com/gookit/color"
)

var prefixError = color.FgRed.Render("Error :")
var prefixWarn = color.FgYellow.Render("Warn :")
var prefixInfo = color.FgGreen.Render("Info :")

// BlueText returns string with blur forefround color
func BlueText(msg ...string) string {
	return color.FgLightBlue.Render(strings.Join(msg, ""))
}

// Fatal logs message with red colored prefix, if `err != nil`, then exits the application
func Fatal(err error, msg ...string) {
	if err != nil {
		log.Fatalf("%s %s - %v\n", prefixError, strings.Join(msg, " "), err)
	}
}

// Error logs message with red colored prefix if `err != nil`
func Error(err error, msg ...string) {
	if err != nil {
		log.Printf("%s %s - %v\n", prefixError, strings.Join(msg, " "), err)
	}
}

// Warn logs message with yellow colored prefix
func Warn(msg ...string) {
	log.Printf("%s %s\n", prefixWarn, strings.Join(msg, " "))
}

// Info logs message with green colored prefix
func Info(msg ...string) {
	log.Printf("%s %s\n", prefixInfo, strings.Join(msg, " "))
}
