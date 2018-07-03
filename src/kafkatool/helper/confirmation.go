package helper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirmation asks the user if he really intends to carry out the concerned action.
func Confirmation(s string) bool {

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		Check(err)

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" {
			return true
		} else if response == "n" {
			return false
		}
	}
}
