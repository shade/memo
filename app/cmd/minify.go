package cmd

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/res"
	"github.com/spf13/cobra"
)

var minifyCmd = &cobra.Command{
	Use:   "minify",
	RunE: func(c *cobra.Command, args []string) error {
		err := res.Minify()
		if err != nil {
			jerr.Get("error minifying js", err).Print()
		}
		return nil
	},
}
