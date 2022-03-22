package backend

import (
	"os"
	"walletboot/chainmgr"
	"walletboot/dao"
	"walletboot/initialize"
	"walletboot/serve"

	"github.com/sirupsen/logrus"
)

type Backend struct {
	Wallet    *serve.Wallet
	Transfer  *serve.Transfer
	XFSClient chainmgr.ChainMgr
}

func NewBackend() *Backend {
	accountDB := initialize.AccountDB()
	transferDB := initialize.TransferDB()
	daokeyDB := dao.NewKeyStoreDB(accountDB)
	wallet, err := serve.NewWallet(daokeyDB)
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}
	// if err := wallet.SetupTxFrom(); err != nil {
	// 	logrus.Warn(err)
	// 	os.Exit(1)
	// }

	transfer := serve.NewTransfer(transferDB)
	return &Backend{
		Wallet:    wallet,
		Transfer:  transfer,
		XFSClient: chainmgr.NewChainMgr(daokeyDB),
	}
}
