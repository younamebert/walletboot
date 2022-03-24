package config

import (
	"math/big"
	"walletboot/common"

	"github.com/shopspring/decimal"
)

var (
	Version                = "0"
	DbPath                 = "./Db"
	LoadAccountsDbPath     = DbPath + "/LoadAccountsDb"
	TxDbPath               = DbPath + "/txDb"
	RpcClientApiHost       = "https://api.scan.xfs.tech/jsonrpc/v2" //"http://127.0.0.1:9012/" 本地的
	RpcClientApiTimeOut    = "180s"
	AccountMaxNumber       = 100
	TxLogPrefix            = []byte("txlog:")
	CronSpec               = "10s"                   // 5s
	AccountFactor          = decimal.NewFromInt(200) // 2%
	BlockTxPoolMaxSize     = int64(100)
	BlockTxPoolMaxSizeShow = true // 打开交易size极限条件
	NewAccountNumber       = 1
	SendTxNumber           = 2
)

var (
	TxGas      = big.NewInt(25000)
	TxGasPrice = big.NewInt(10)
)

func DefaultGasPrice() *big.Int {
	return common.NanoCoin2Atto(TxGasPrice)
}
