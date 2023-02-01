package cmd

import (
	"github.com/spf13/cobra"
)

var (
	postCmd = &cobra.Command{
		Use:   "post",
		Short: "Make a post",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}
)
