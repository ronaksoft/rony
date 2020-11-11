package cmd

import (
	"github.com/ronaksoft/rony/config"
	"github.com/spf13/cobra"
)

/*
   Creation Time: 2020 - Aug - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func init() {
	config.SetPersistentFlags(RootCmd,
		config.BoolFlag("dryRun", false, "if set then it is a dry run for testing purpose"),
	)
}

var RootCmd = &cobra.Command{
	Use: "rony",
}
