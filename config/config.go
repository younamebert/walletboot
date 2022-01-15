package config

import "math/big"

var (
	DbPath              = "./Db"
	LoadAccountsDbPath  = DbPath + "/LoadAccountsDb"
	TxDbPath            = DbPath + "/txDb"
	RpcClientApiHost    = "http://127.0.0.1:9012/"
	RpcClientApiTimeOut = "180s"
	CronSpec            = "7s"                                     // 5s
	AccountFactor       = new(big.Float).SetFloat64(float64(0.02)) // 2%
	NewAccountNumber    = 1                                        //
	SendTxNumber        = 4                                        //
)
