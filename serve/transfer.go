package serve

import (
	"encoding/json"
	"errors"
	"strconv"
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
	Conn          *httpxfs.Client
	TargetHeight  string
	CurrentHeight string
	CurrentHash   string
}

type Txlog struct {
	From          string
	To            string
	Value         string
	TxHash        string
	TargetHeight  string
	CurrentHeight string
	CurrentHash   string
}

func NewTxSend(txDb badger.IStorage, cli *httpxfs.Client) *Transfer {
	return &Transfer{
		TxDb: txDb,
		Conn: cli,
	}
}

//上一个区块交易成功的区块
func (t *Transfer) GetCurrentTxs() ([]map[string]string, error) {
	// 获取当前主链最高块
	block := make(map[string]interface{})
	if err := t.Conn.CallMethod(1, "Chain.Head", nil, &block); err != nil {
		return nil, err
	}
	if block == nil {
		return nil, errors.New("height does not exist")
	}
	// 更新交易日志信息
	blockHeight := block["height"].(string)
	blockNumber, err := strconv.ParseInt(blockHeight, 10, 64)
	if err != nil {
		return nil, err
	}

	TargetNumber := blockNumber + 1
	t.CurrentHeight = blockHeight
	t.CurrentHash = block["hash"].(string)
	t.TargetHeight = strconv.FormatInt(TargetNumber, 10)

	// 获取最高块交易成功的数据
	type getTxsByBlockNumArgs struct {
		Number string `json:"number"`
	}

	req := &getTxsByBlockNumArgs{
		Number: blockHeight,
	}

	result := make([]map[string]interface{}, 0)
	if err := t.Conn.CallMethod(1, "Chain.GetTxsByBlockNum", &req, &result); err != nil {
		return nil, err
	}

	resp := make([]map[string]string, 0)
	for _, v := range result {
		deposit := make(map[string]string, 0)
		to := v["to"].(string)
		deposit[to] = v["value"].(string)
		resp = append(resp, deposit)
	}
	return resp, nil
}

func (t *Transfer) SendTransactionFunc(args *SendTransactionArgs) (string, error) {

	var txhash *string
	if err := t.Conn.CallMethod(1, "Wallet.SendTransaction", args, &txhash); err != nil {
		return "", err
	}

	if err := t.writeTxLog(txhash, args); err != nil {
		return "", err
	}
	return *txhash, nil
}

func (t *Transfer) writeTxLog(txhash *string, args *SendTransactionArgs) error {

	txlog := &Txlog{
		From:          args.From,
		To:            args.To,
		Value:         args.Value,
		TxHash:        *txhash,
		TargetHeight:  t.TargetHeight,
		CurrentHeight: t.CurrentHeight,
		CurrentHash:   t.CurrentHash,
	}
	bs, err := json.Marshal(txlog)
	if err != nil {
		return err
	}
	return t.TxDb.Set(args.From, bs)
}
