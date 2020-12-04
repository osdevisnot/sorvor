// Package logger provides color coded log messages with error handling
package logger

import (
	"log"
	"strings"

	"github.com/gookit/color"
)

var prefixError = color.FgRed.Render("error :")
var prefixWarn = color.FgYellow.Render("warn :")
var prefixInfo = color.FgGreen.Render("info :")

// Fatal logs message with red colored prefix and exits the program if `err != nil`
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

// BlueText returns string with blur foreground color
func BlueText(msg ...string) string {
	return color.FgLightBlue.Render(strings.Join(msg, ""))
}
