package utils

import (
	"fmt"
	"log"
	"strings"
)

func CheckError(e error, message string) {
	if e != nil {
		log.Fatalln(message, ":", e)
	}
}

func MakeFirstUpper(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func FormatValue(value float64) string {
	var suffix string
	var divisor float64

	switch {
	case value >= 1e6:
		suffix = "M"
		divisor = 1e6
	case value >= 1e3:
		suffix = "K"
		divisor = 1e3
	default:
		return fmt.Sprintf("%f", value)
	}

	return fmt.Sprintf("%.1f%s", value/divisor, suffix)
}
