package main

import (
	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/utils/osutil"
)

func main() {
	cobra.SetMousetrapHelpText("")

	osutil.ExitErr(rootCmd.Execute())
}
