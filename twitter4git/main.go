package main

import (
	"fmt"
	"os"

	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/must"
	"github.com/petar/twitter4git/twitter4git/cmd"
)

func main() {
	if base.IsVerbose() {
		cmd.Execute()
	} else {
		err := must.Try(
			func() { cmd.Execute() },
		)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
}
