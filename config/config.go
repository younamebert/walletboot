package config

import (
	"math/big"
	"math/rand"
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
	CronSpec               = "20s"                   // 5s
	AccountFactor          = decimal.NewFromInt(200) // 2%
	BlockTxPoolMaxSize     = int64(30)
	BlockTxPoolMaxSizeShow = true // 打开交易size极限条件
	NewAccountNumber       = 0
	SendTxNumber           = 1
)

var (
	TxGas      = big.NewInt(25000)
	TxGasPrice = big.NewInt(randGasPrice())
)

func randGasPrice() int64 {
	return int64(rand.Intn(20) + 10)
}
func DefaultGasPrice() *big.Int {
	return common.NanoCoin2Atto(TxGasPrice)
}
