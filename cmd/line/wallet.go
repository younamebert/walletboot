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

// 	walletGetAddrPriKeyCommand = &cobra.Command{
// 		Use:                   "addrprikey <address>",
// 		DisableFlagsInUseLine: true,
// 		Short:                 "get account address prikey ",
// 		RunE:                  getAddrPriKey,
// 	}
// )

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

// func getAddrPriKey(cmd *cobra.Command, args []string) error {
// 	// cli := client.NewClient(config.RpcClientApiHost, config.RpcClientApiTimeOut)
// 	// cli.CallMethod()
// 	return nil
// }

// func walletNumber() error {
// 	number := task.AppCore().Wallet.GetNumber()
// 	fmt.Printf("%v\n", number)
// 	return nil
// }

// func init() {
// 	// walletCommand.AddCommand(walletNumberCommand)
// 	walletCommand.AddCommand(walletGetAddrPriKeyCommand)
// 	rootCmd.AddCommand(walletCommand)
// }
