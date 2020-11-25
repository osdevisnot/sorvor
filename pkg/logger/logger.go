package logger

import (
	"github.com/gookit/color"
	"log"
	"strings"
)

var prefixError = color.FgRed.Render("Error :")
var prefixWarn = color.FgYellow.Render("Warn :")
var prefixInfo = color.FgGreen.Render("Info :")

func BlueText(msg ...string) string {
	return color.FgLightBlue.Render(strings.Join(msg, ""))
}

func Fatal(err error, msg ...string) {
	if err != nil {
		log.Fatalf("%s %s - %v\n", prefixError, strings.Join(msg, " "), err)
	}
}

func Error(err error, msg ...string) {
	if err != nil {
		log.Printf("%s %s - %v\n", prefixError, strings.Join(msg, " "), err)
	}
}

func Warn(msg ...string) {
	log.Printf("%s %s\n", prefixWarn, strings.Join(msg, " "))
}

func Info(msg ...string) {
	log.Printf("%s %s\n", prefixInfo, strings.Join(msg, " "))
}
