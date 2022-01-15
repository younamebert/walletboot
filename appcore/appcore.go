package appcore

import (
	"math/big"
	"walletboot/common"
	"walletboot/config"
	"walletboot/httpxfs"
	"walletboot/serve"
	"walletboot/storage/badger"

	"github.com/sirupsen/logrus"
)

type AppCore struct {
	Wallet   *serve.Wallet
	Transfer *serve.Transfer
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

	wallet := serve.NewWallet(LoadAccountsDb)

	transfer := serve.NewTxSend(txDb, cli)

	return &AppCore{
		Wallet:   wallet,
		Transfer: transfer,
	}, nil
}

// 生成钱包
func (c *AppCore) RunRand() {
	if _, err := c.Wallet.NewAccount(); err != nil {
		logrus.Error(err)
		return
	}
}

func (c *AppCore) CoreWalletFrom() string {
	addr, err := c.Wallet.NewAccount()
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return addr.B58String()
}

// 发送交易
func (c *AppCore) RunSendTx() {
	c.UpdateAccount()
	txTo, txFromObj := c.Wallet.RandAddr()
	request := &serve.SendTransactionArgs{
		To: txTo,
	}
	for addr, val := range txFromObj {
		request.From = addr
		request.Value = c.randAmount(val)
	}

	hash, err := c.Transfer.SendTransactionFunc(request)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("send Tx From:%v To:%v value:%v txHash:%v", request.From, request.To, request.Value, hash)
}

func (c *AppCore) UpdateAccount() {
	list, err := c.Transfer.GetCurrentTxs()
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, v := range list {
		for ks, vs := range v {
			addr := common.B58ToAddress([]byte(ks))
			if err := c.Wallet.UpdateAccout(addr, vs); err != nil {
				logrus.Error(err)
				continue
			}
		}
	}
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
