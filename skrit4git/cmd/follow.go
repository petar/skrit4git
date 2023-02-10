package cmd

import (
	"github.com/petar/skrit4git/proto"
	"github.com/spf13/cobra"
)

var (
	followCmd = &cobra.Command{
		Use:   "follow",
		Short: "Follow a user",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			proto.Follow(ctx, setup.Home, proto.MustParseHandle(ctx, followHandle))
		},
	}
)

var (
	followHandle string
)

func init() {
	rootCmd.AddCommand(followCmd)
	followCmd.Flags().StringVarP(&followHandle, "handle", "h", "", "user handle to follow (e.g. maymounkov.org)")
	followCmd.MarkFlagRequired("handle")
}
