package chainmgr

import (
	"walletboot/chainmgr/client"
	"walletboot/config"
	"walletboot/serve"

	"github.com/sirupsen/logrus"
)

type ChainMgr interface {
	CurrentBHeader() *BlockHeader
	GetReceiptByHash(txhash string) *Receipt
	GetNonce(addr string) int64
	SendRawTransaction(data string) *string
	GetAccount(addr string) *serve.Accounts
	// GetTxsByBlockHash(blockHash string) []*Transaction
	// GetBlockHeaderByHash(blockHash string) *BlockHeader
	// GetAccountInfo(address string) *AccountState
	// GetBlockHeaderByNumber(blocknumber string) *BlockHeader
}

type ChainMgrs struct {
	xfsClient *client.Client
}

func NewChainMgr() *ChainMgrs {
	cli := client.NewClient(config.RpcClientApiHost, config.RpcClientApiTimeOut)
	return &ChainMgrs{
		xfsClient: cli,
	}
}

func (ext *ChainMgrs) CurrentBHeader() *BlockHeader {
	lastBlockHeader := new(BlockHeader)
	if err := ext.xfsClient.CallMethod(1, "Chain.Head", nil, &lastBlockHeader); err != nil {
		return nil
	}
	if lastBlockHeader != nil {
		return lastBlockHeader
	}
	return nil
}

func (ext *ChainMgrs) GetReceiptByHash(txhash string) *Receipt {
	req := &GetTxByHashArgs{
		TxHash: txhash,
	}
	recs := new(Receipt)
	if err := ext.xfsClient.CallMethod(1, "Chain.GetReceiptByHash", &req, &recs); err != nil {
		logrus.Warn(err)
		return nil
	}
	if recs != nil {
		return recs
	}
	return nil
}

func (ext *ChainMgrs) GetNonce(addr string) int64 {
	req := &GetAddrNonceByHashArgs{
		Address: addr,
	}
	var result int64 = 0
	if err := ext.xfsClient.CallMethod(1, "TxPool.GetAddrTxNonce", &req, &result); err != nil {
		logrus.Warn(err)
		return result
	}
	return result
}

func (ext *ChainMgrs) SendRawTransaction(data string) *string {
	var txhash *string
	args := &SendRawTxArgs{
		Data: data,
	}
	if err := ext.xfsClient.CallMethod(1, "TxPool.SendRawTransaction", args, &txhash); err != nil {
		logrus.Warn(err)
		return nil
	}
	return txhash
}

func (ext *ChainMgrs) GetAccount(addr string) *serve.Accounts {
	req := &GetAccountArgs{
		Address: addr,
	}

	chainStatusLast := &serve.Accounts{}
	if err := ext.xfsClient.CallMethod(1, "State.GetAccount", &req, &chainStatusLast); err != nil {
		return nil
	}
	return chainStatusLast
}

// var txhash string

// 	if err := t.updateWriteMsg(); err != nil {
// 		return "", err
// 	}

// 	if err := t.conn.CallMethod(1, "TxPool.SendRawTransaction", args, &txhash); err != nil {
// 		return "", err
// 	}
