package appcore

import (
	"math/big"
	"testtx/common"
	"testtx/config"
	"testtx/httpxfs"
	"testtx/serve"
	"testtx/storage/badger"

	"github.com/sirupsen/logrus"
)

type AppCore struct {
	Wallet *serve.RandWallet
	SendTx *serve.Txlog
}

func New() (*AppCore, error) {
	LoadAccountsDb, err := badger.New(config.LoadAccountsDbPath)
	if err != nil {
		return nil, err
	}

	txDb, err := badger.New(config.TxDbPath)
	if err != nil {
		return nil, err
	}

	cli := httpxfs.NewClient(config.RpcClientApiHost, config.RpcClientApiTimeOut)

	wallet := serve.NewRqWallet(LoadAccountsDb, cli)

	txSend := serve.NewTxSend(txDb, cli)

	return &AppCore{
		Wallet: wallet,
		SendTx: txSend,
	}, nil
}

//new Wallet
func (c *AppCore) RunRand() {
	if err := c.Wallet.RandCreateWallet(); err != nil {
		logrus.Error(err)
		return
	}
}

// send transfer
func (c *AppCore) RunSendTx() {
	request := &serve.SendTransactionArgs{
		To: c.Wallet.GetTxTo(),
	}
	for k, v := range c.Wallet.GetTxFrom() {
		request.From = k
		request.Value = c.randAmount(v)
	}
	hash, err := c.SendTx.SendTransactionFunc(request)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("send Tx From:%v To:%v value:%v txHash:%v", request.From, request.To, request.Value, hash)
}

func (c *AppCore) randAmount(val string) string {
	result := big.NewFloat(0)

	bal, err := common.Atto2BaseRatCoin(val)
	if err != nil {
		return "0"
	}
	result = result.Mul(bal, config.AccountFactor)

	return result.String()

}
