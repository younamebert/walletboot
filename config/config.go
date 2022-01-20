package config

import (
	"math/big"
	"walletboot/common"
)

var (
	Version             = "0"
	DbPath              = "./Db"
	LoadAccountsDbPath  = DbPath + "/LoadAccountsDb"
	TxDbPath            = DbPath + "/txDb"
	RpcClientApiHost    = "http://127.0.0.1:9012/"
	RpcClientApiTimeOut = "180s"
	AccountMaxNumber    = 100
	TxLogPrefix         = []byte("txlog:")
	CronSpec            = "3s"                            // 5s
	AccountFactor       = new(big.Float).SetFloat64(0.02) // 2%
	NewAccountNumber    = 1                               //
	SendTxNumber        = 1                               //
)

var (
	TxGas      = big.NewInt(25000)
	TxGasPrice = big.NewInt(10)
)

func DefaultGasPrice() *big.Int {
	return common.NanoCoin2Atto(TxGasPrice)
}
