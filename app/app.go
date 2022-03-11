package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"walletboot/backend"
	"walletboot/chainmgr"
	"walletboot/common"
	"walletboot/config"
	"walletboot/serve"

	"github.com/sirupsen/logrus"
)

type App struct {
	back      *backend.Backend
	xfsClient chainmgr.ChainMgr
}

// var Cil *httpxfs.Client

func New() *App {
	return &App{
		back:      backend.NewBackend(),
		xfsClient: chainmgr.NewChainMgr(),
	}
}

// Create a new account
func (app *App) CreateAccount() error {
	if _, err := app.back.Wallet.CreateAccount(); err != nil {
		return err
	}
	return nil
}

//Send transaction
func (app *App) SendTransaction() error {

	accounts := app.back.Wallet.RandomAccessAccount()

	if len(accounts) == 0 {
		return fmt.Errorf("no user with from qualification was found in the wallet launcher")
	}

	tx, req, err := app.NewTransaction(accounts[0], accounts[1])
	if err != nil {
		return err
	}
	hash := app.xfsClient.SendRawTransaction(req.Data)
	if hash != nil {
		if err := app.back.Transfer.WriteTxLog(*hash, tx); err != nil {
			return err
		}
		logrus.Infof("send Tx From:%v To:%v value:%v txHash:%v", accounts[1].Address, tx.To, tx.Value, hash)
		return nil
	}
	return nil
}

func (app *App) NewTransaction(toObj, fromObj *serve.Accounts) (*serve.SendTransaction, *serve.SendRawTxArgs, error) {
	tx := &serve.SendTransaction{}
	req := &serve.SendRawTxArgs{}
	priKey, err := app.back.Wallet.GetKeyByAddress(common.B58ToAddress([]byte(fromObj.Address)))
	if err != nil {
		return nil, nil, err
	}

	nonce := app.xfsClient.GetNonce(fromObj.Address)

	tx.Version = config.Version
	tx.To = toObj.Address
	tx.Nonce = strconv.FormatInt(nonce, 10)
	tx.GasLimit = config.TxGas.String()
	tx.GasPrice = config.DefaultGasPrice().String()
	tx.Value = app.randAmount(fromObj.Balance.String())
	sign, err := app.back.Transfer.SignHash(tx, priKey)
	if err != nil {
		return nil, nil, err
	}
	tx.Signature = sign
	bs, err := json.Marshal(tx)
	if err != nil {
		return nil, nil, err
	}
	req.Data = base64.StdEncoding.EncodeToString(bs)

	return tx, req, nil
}

func (app *App) randAmount(val string) string {

	result := big.NewFloat(0)

	bal, err := common.Atto2BaseRatCoin(val)
	if err != nil {
		return "0"
	}
	result = result.Mul(bal, config.AccountFactor)

	return result.String()
}