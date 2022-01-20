package appcore

import (
	"encoding/base64"
	"encoding/json"
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

	if len(txFromObj) == 0 {
		return fmt.Errorf("no user with from qualification was found in the wallet launcher")
	}

	tx := &serve.SendTransaction{}
	req := &serve.SendRawTxArgs{}
	for addr, val := range txFromObj {

		priKey, err := c.Wallet.GetKeyByAddress(common.B58ToAddress([]byte(addr)))
		if err != nil {
			return err
		}

		nonce, err := c.Transfer.GetNonce(addr)
		if err != nil {
			return err
		}
		tx.Version = config.Version
		tx.To = txTo
		tx.Nonce = nonce
		tx.GasLimit = config.TxGas.String()
		tx.GasPrice = config.DefaultGasPrice().String()
		tx.Value = c.randAmount(val)
		sign, err := c.Transfer.SignHash(tx, priKey)
		if err != nil {
			return err
		}
		tx.Signature = sign
		bs, err := json.Marshal(tx)
		if err != nil {
			return err
		}
		req.Data = base64.StdEncoding.EncodeToString(bs)

		hash, err := c.Transfer.SendTransactionFunc(req)
		if err != nil {
			return err
		}

		if err := c.Transfer.WriteTxLog(addr, hash, tx); err != nil {
			return err
		}
		logrus.Infof("send Tx From:%v To:%v value:%v txHash:%v", addr, tx.To, tx.Value, hash)
	}
	return nil
}

func (c *AppCore) randAmount(val string) string {
	remain := big.NewFloat(0)

	// bal, ok := new(big.Int).SetString(val, 0)
	bal, err := common.Atto2BaseRatCoin(val)
	if err != nil {
		return "0"
	}
	remain = remain.Mul(bal, config.AccountFactor)

	result, err := common.BaseCoin2Atto(remain.String())
	if err != nil {
		return "0"
	}
	// remain = remain.(bal, config.AccountFactor)
	return result.String()
}
