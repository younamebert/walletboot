package serve

import (
	"encoding/json"
	"fmt"
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

type Transfer struct {
	TxDb          badger.IStorage
	conn          *httpxfs.Client
	CurrentHeight string
	CurrentHash   string
}

type Txlog struct {
	From          string
	To            string
	Value         string
	TxHash        string
	CurrentHeight string
	CurrentHash   string
}

func NewTxSend(txDb badger.IStorage, cli *httpxfs.Client) *Transfer {
	return &Transfer{
		TxDb: txDb,
		conn: cli,
	}
}

func (t *Transfer) SendTransactionFunc(args *SendTransactionArgs) (string, error) {

	if err := t.updateWriteMsg(); err != nil {
		return "", err
	}

	var txhash string

	if err := t.conn.CallMethod(1, "Wallet.SendTransaction", args, &txhash); err != nil {
		return "", err
	}

	if err := t.writeTxLog(txhash, args); err != nil {
		return "", err
	}
	return txhash, nil
}

func (t *Transfer) updateWriteMsg() error {
	head := make(map[string]interface{})
	err := t.conn.CallMethod(1, "Chain.Head", nil, &head)
	if err != nil {
		return err
	}

	t.CurrentHash = head["hash"].(string)
	t.CurrentHeight = fmt.Sprint(head["height"])
	return nil
}

func (t *Transfer) writeTxLog(txhash string, args *SendTransactionArgs) error {

	txlog := &Txlog{
		From:          args.From,
		To:            args.To,
		Value:         args.Value,
		TxHash:        txhash,
		CurrentHeight: t.CurrentHeight,
		CurrentHash:   t.CurrentHash,
	}
	bs, err := json.Marshal(txlog)
	if err != nil {
		return err
	}
	return t.TxDb.Set(args.From, bs)
}
