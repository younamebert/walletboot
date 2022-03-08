package serve

import (
	"encoding/json"
	"testtx/httpxfs"
	"testtx/storage/badger"
)

type SendTransactionArgs struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type GetAccountArgs struct {
	RootHash string `json:"root_hash"`
	Address  string `json:"address"`
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
