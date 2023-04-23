package cmd

import "fmt"

func handleCmdErr(err error) {
	if err != nil {
		if err.Error() == "^D" || err.Error() == "^C" {
			fmt.Println("GoodBye")
			return
		}
		Log.Errorf("Error: %v", err)
	}
}
