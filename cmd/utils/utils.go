package utils

import "log"

func CheckError(e error, message string) {
	if e != nil {
		log.Fatalln(message, ":", e)
	}
}
