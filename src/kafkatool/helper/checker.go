// Package helper implements general helper methods.
package helper

import "log"

// Check checks the error
func Check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
