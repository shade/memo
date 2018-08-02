package cmd

import (
	"fmt"
	"github.com/memocash/memo/app/bitcoin/main-node"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var actionNodeCmd = &cobra.Command{
	Use: "action-node",
	RunE: func(c *cobra.Command, args []string) error {
		var last time.Time
		for last.IsZero() || time.Since(last) > time.Minute {
			last = time.Now()
			main_node.StartActionNode()
			main_node.WaitForActionNodeDisconnect()
			fmt.Println("Disconnected.")
			main_node.ActionNode = main_node.Node{}
		}
		os.Exit(1)
		return nil
	},
}

var userNodeCmd = &cobra.Command{
	Use: "user-node",
	RunE: func(c *cobra.Command, args []string) error {
		var last time.Time
		for last.IsZero() || time.Since(last) > time.Minute {
			last = time.Now()
			main_node.StartUserNode()
			main_node.WaitForUserNodeDisconnect()
			fmt.Println("Disconnected.")
			main_node.UserNode = main_node.Node{}
		}
		os.Exit(1)
		return nil
	},
}
