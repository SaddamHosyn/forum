package util

import "log"

// checkErr checks if an error occurred and logs the provided message if it did.
func CheckErr(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
