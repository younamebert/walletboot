package serve

import (
	"encoding/json"
	"fmt"
	"walletboot/config"
	"walletboot/httpxfs"
	"walletboot/storage/badger"
)

type SendTransactionArgs struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type GetAccountArgs struct {
	Address string `json:"address"`
}

type Txlog struct {
	TxDb       badger.IStorage
	ClientHttp *httpxfs.Client
}

func NewTxSend(txDb badger.IStorage, cli *httpxfs.Client) *Txlog {
	return &Txlog{
		TxDb:       txDb,
		ClientHttp: cli,
	}
}

func (t *Txlog) GetTxPoolPending() error {
	var txPoolPendingCount int
	if err := t.ClientHttp.CallMethod(1, "TxPool.GetPendingSize", nil, &txPoolPendingCount); err != nil {
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	if txPoolPendingCount >= config.TxPoolPendingNumberMax {
		return fmt.Errorf("pending max number:%v thisnumber:%v\n", config.TxPoolPendingNumberMax, txPoolPendingCount)
	}
	return nil
}
func (t *Txlog) SendTransactionFunc(args *SendTransactionArgs) (string, error) {

	var result *string
	if err := t.ClientHttp.CallMethod(1, "Wallet.SendTransaction", args, &result); err != nil {
		return "", err
	}
	bs, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	if err := t.TxDb.Set(*result, bs); err != nil {
		return "", err
	}
	return *result, nil
}
