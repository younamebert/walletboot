package line

// import (
// 	"fmt"
// 	"walletboot/common"

// 	"github.com/spf13/cobra"
// )

// var (
// 	walletCommand = &cobra.Command{
// 		Use:                   "wallet [options]",
// 		DisableFlagsInUseLine: true,
// 		SilenceUsage:          true,
// 		Short:                 "get wallet info",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			return cmd.Help()
// 		},
// 	}
// 	walletFormsCommand = &cobra.Command{
// 		Use:                   "list [options]",
// 		DisableFlagsInUseLine: true,
// 		Short:                 "get wallet address list",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			return walletForms()
// 		},
// 	}
// 	walletNumberCommand = &cobra.Command{
// 		Use:                   "accountNumber [options]",
// 		DisableFlagsInUseLine: true,
// 		Short:                 "get wallet address number",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			return walletNumber()
// 		},
// 	}
// )

// func walletForms() error {
// 	result, err := task.AppCore().Wallet.GetForms()
// 	if err != nil {
// 		return err
// 	}
// 	bs, err := common.MarshalIndent(result)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(string(bs))
// 	return nil
// }

// func walletNumber() error {
// 	number := task.AppCore().Wallet.GetNumber()
// 	fmt.Printf("%v\n", number)
// 	return nil
// }

// func init() {
// 	walletCommand.AddCommand(walletNumberCommand)
// 	walletCommand.AddCommand(walletFormsCommand)
// 	rootCmd.AddCommand(walletCommand)
// }
