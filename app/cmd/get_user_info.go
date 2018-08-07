package cmd

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
	"github.com/spf13/cobra"
)

var getUserInfoCmd = &cobra.Command{
	Use: "info",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) != 1 {
			return jerr.New("invalid number of arguments, must give a username")
		}
		var username = args[0]
		user, err := db.GetUserByUsername(username)
		if err != nil {
			jerr.Get("error getting user by username", err).Print()
			return nil
		}
		key, err := db.GetKeyForUser(user.Id)
		if err != nil {
			jerr.Get("error getting key for user", err).Print()
			return nil
		}
		fmt.Printf("User: %s (id: %d)\nAddress: %s\n", user.Username, user.Id, key.GetAddress().GetEncoded())
		return nil
	},
}
