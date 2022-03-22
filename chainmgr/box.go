package chainmgr

import (
	"encoding/json"
	"fmt"
	"walletboot/chainmgr/client"
	"walletboot/common"
	"walletboot/config"
	"walletboot/dao"
	"walletboot/serve"

	"github.com/sirupsen/logrus"
)

type ChainMgr interface {
	CurrentBHeader() *BlockHeader
	GetReceiptByHash(txhash string) *Receipt
	GetNonce(addr string) int64
	SendRawTransaction(data string) *string
	GetAccount(addr string) *serve.Accounts
	UpdateAccountState() error
	// GetTxsByBlockHash(blockHash string) []*Transaction
	// GetBlockHeaderByHash(blockHash string) *BlockHeader
	// GetAccountInfo(address string) *AccountState
	// GetBlockHeaderByNumber(blocknumber string) *BlockHeader
}

type ChainMgrs struct {
	xfsClient *client.Client
	db        *dao.KeyStoreDB
}

func NewChainMgr(daokeyDB *dao.KeyStoreDB) *ChainMgrs {
	cli := client.NewClient(config.RpcClientApiHost, config.RpcClientApiTimeOut)
	return &ChainMgrs{
		xfsClient: cli,
		db:        daokeyDB,
	}
}

func (ext *ChainMgrs) UpdateAccountState() error {
	return ext.db.AddrForeach(func(k string, v []byte) error {
		// Loop to update the status of all users
		accounts := &serve.Accounts{}

		if err := json.Unmarshal(v, &accounts); err != nil {
			return err
		}

		req := &getAccountArgs{
			Address: accounts.Address,
		}

		chainStatusLast := &serve.Accounts{}
		if err := ext.xfsClient.CallMethod(1, "State.GetAccount", &req, &chainStatusLast); err != nil {
			return err
		}
		accounts.Nonce = chainStatusLast.Nonce
		accounts.Extra = chainStatusLast.Extra
		accounts.StateRoot = chainStatusLast.StateRoot
		accounts.Code = chainStatusLast.Code
		accounts.Balance = chainStatusLast.Balance

		addr := common.B58ToAddress([]byte(accounts.Address))

		_, err := ext.db.GetAccount(addr)
		if err != nil {
			return err
		}
		//
		bs, err := json.Marshal(accounts)
		if err != nil {
			return err
		}

		if err := ext.db.UpdateAccount(addr, bs); err != nil {
			return err
		}
		return nil
	})
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
	fmt.Println(args.Data)
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
