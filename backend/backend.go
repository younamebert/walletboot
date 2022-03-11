package backend

import (
	"os"
	"walletboot/initialize"
	"walletboot/serve"

	"github.com/sirupsen/logrus"
)

type Backend struct {
	Wallet   *serve.Wallet
	Transfer *serve.Transfer
}

func NewBackend() *Backend {

	accountDB := initialize.AccountDB()
	transferDB := initialize.TransferDB()
	// cli := httpxfs.NewClient(config.RpcClientApiHost, config.RpcClientApiTimeOut)

	wallet, err := serve.NewWallet(accountDB)
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}
	if err := wallet.SetupTxFrom(); err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}

	transfer := serve.NewTransfer(transferDB)
	return &Backend{
		Wallet:   wallet,
		Transfer: transfer,
	}
}
