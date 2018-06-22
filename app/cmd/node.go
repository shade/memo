package cmd

import (
	"fmt"
	"github.com/memocash/memo/app/bitcoin/main-node"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var mainNodeCmd = &cobra.Command{
	Use: "main-node",
	RunE: func(c *cobra.Command, args []string) error {
		var last time.Time
		for last.IsZero() || time.Since(last) > time.Minute {
			last = time.Now()
			main_node.Start()
			main_node.WaitForDisconnect()
			fmt.Println("Disconnected.")
			main_node.BitcoinNode = main_node.Node{}
		}
		os.Exit(1)
		return nil
	},
}
