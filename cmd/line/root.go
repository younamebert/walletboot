package line

import (
	"fmt"
	"os"
	"walletboot/bootcron"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:                   fmt.Sprintf("%s <command> [<options>]", "walletboot"),
		DisableFlagsInUseLine: true,
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

var (
	task, err = bootcron.New()
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		os.Exit(1)
	}
}
