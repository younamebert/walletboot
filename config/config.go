package config

import "math/big"

var (
	LoadAccountsDbPath     = "./LoadAccountsDb"
	TxDbPath               = "./txDb"
	RpcClientApiHost       = "http://127.0.0.1:9012/"
	RpcClientApiTimeOut    = "180s"
	CronSpec               = "20s"                                    //每隔5秒执行一次
	AccountFactor          = new(big.Float).SetFloat64(float64(0.02)) // 2%
	AccountNumberMax       = 30
	TxPoolPendingNumberMax = 50
	NewAccountNumber       = 1 // 批量 生成钱包
	SendTxNumber           = 1 //批量 发送交易数据
)
