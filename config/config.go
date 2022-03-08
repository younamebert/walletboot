package config

import "math/big"

var (
	DbPath              = "./Db"
	LoadAccountsDbPath  = DbPath + "/LoadAccountsDb"
	TxDbPath            = DbPath + "/txDb"
	RpcClientApiHost    = "http://127.0.0.1:9012/"
	RpcClientApiTimeOut = "180s"
	NewAccountNumberMax = 100
	CronSpec            = "3s"                                     //每隔5秒执行一次
	AccountFactor       = new(big.Float).SetFloat64(float64(0.02)) // 2%
	NewAccountNumber    = 0                                        // 批量 生成钱包
	SendTxNumber        = 1                                        //批量 发送交易数据
)
