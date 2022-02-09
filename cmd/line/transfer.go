package line

import (
	"fmt"
	"walletboot/common"

	"github.com/spf13/cobra"
)

var (
	daemonCmd = &cobra.Command{
		Use:                   "walletboot [options]",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Short:                 "Start a send transfer process",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	transferBootStartCmd = &cobra.Command{
		Use:                   "Start",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Short:                 "Start a send transfer process",
		RunE: func(cmd *cobra.Command, args []string) error {
			return start()
		},
	}
	transferBootStopCmd = &cobra.Command{
		Use:                   "Stop",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Short:                 "Stop a send transfer process",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stop()
		},
	}
	txlogListCmd = &cobra.Command{
		Use:                   "Stop",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Short:                 "Stop a send transfer process",
		RunE: func(cmd *cobra.Command, args []string) error {
			return txlogList()
		},
	}
)

// process
func start() error {
	if err != nil {
		return err
	}
	go task.Start()
	select {}
}

func stop() error {
	task.Stop()
	return nil
}

func txlogList() error {
	txlogjson := task.AppCore().Transfer.ListTxLog()
	bs, err := common.MarshalIndent(txlogjson)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", string(bs))
	return nil
}

func init() {
	// mFlags := daemonCmd.PersistentFlags()
	daemonCmd.AddCommand(txlogListCmd)
	daemonCmd.AddCommand(transferBootStartCmd)
	daemonCmd.AddCommand(transferBootStopCmd)
	rootCmd.AddCommand(daemonCmd)
}
