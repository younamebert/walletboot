package initialize

import (
	"os"
	"walletboot/config"
	"walletboot/storage/badger"

	"github.com/sirupsen/logrus"
)

func AccountDB() *badger.Storage {
	accountsDB, err := badger.New(config.LoadAccountsDbPath)
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}
	return accountsDB
}

func TransferDB() *badger.Storage {
	transferDB, err := badger.New(config.TxDbPath)
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}
	return transferDB
}
