package cmd

import (
	"fmt"
	"os"

	"github.com/petar/skrit4git/proto"
	"github.com/spf13/cobra"
)

var (
	postCmd = &cobra.Command{
		Use:   "post",
		Short: "Make a post",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			chg := proto.Post(ctx, setup.Home, XXXcontent)
			fmt.Fprint(os.Stdout, setup.Home.Link(chg.Result))
		},
	}
)

var (
	postContent string
)

func init() {
	rootCmd.AddCommand(postCmd)
	followCmd.Flags().StringVarP(&postContent, "content", "c", "", "post content")
}
