package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gov4git/lib4git/must"
	"github.com/petar/skrit4git/proto"
	"github.com/spf13/cobra"
)

var (
	postCmd = &cobra.Command{
		Use:   "post",
		Short: "Make a post",
		Long: `By default, the contents of the post is read from the standard input.
If the --msg flag is specified, its value is the post.
If the --file flag is specific, its value is the name of a file containing the post.`,
		Run: func(cmd *cobra.Command, args []string) {
			var content string
			switch {
			case postMsg != "":
				content = postMsg
			case postFile != "":
				buf, err := ioutil.ReadFile(postFile)
				must.NoError(ctx, err)
				content = string(buf)
			default:
				buf, err := ioutil.ReadAll(os.Stdin)
				must.NoError(ctx, err)
				content = string(buf)
			}
			chg := proto.Post(ctx, setup.Home, content)
			fmt.Fprint(os.Stdout, setup.Home.Link(chg.Result))
		},
	}
)

var (
	postFile string
	postMsg  string
)

func init() {
	rootCmd.AddCommand(postCmd)
	followCmd.Flags().StringVarP(&postMsg, "msg", "m", "", "post message")
	followCmd.Flags().StringVarP(&postFile, "file", "f", "", "file containing post")
}
