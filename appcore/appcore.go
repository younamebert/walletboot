package appcore

import (
	"fmt"
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
	fmt.Println(0)
	wallet, err := serve.NewWallet(LoadAccountsDb, cli)
	if err != nil {
		return nil, err
	}

	transfer := serve.NewTxSend(txDb, cli)
	return &AppCore{
		Wallet:   wallet,
		Transfer: transfer,
	}, nil
}

//
func (c *AppCore) RunRand() error {
	if _, err := c.Wallet.NewAccount(); err != nil {
		return err
	}
	return nil
}

//Send transaction
func (c *AppCore) RunSendTx() error {

	if err := c.Wallet.UpdateAccount(); err != nil {
		return err
	}

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
		return err
	}
	logrus.Infof("send Tx From:%v To:%v value:%v txHash:%v", request.From, request.To, request.Value, hash)
	return nil
}

func (c *AppCore) UpdateAccount() {
	list, err := c.Wallet.GetCurrentTxs()
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
