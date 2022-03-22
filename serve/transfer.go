package serve

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"walletboot/common"
	"walletboot/common/ahash"
	"walletboot/config"
	"walletboot/crypto"
	"walletboot/storage/badger"
)

type SendTransaction struct {
	Version   string `json:"version"`
	To        string `json:"to"`
	Value     string `json:"value"`
	GasPrice  string `json:"gas_price"`
	GasLimit  string `json:"gas_limit"`
	Data      string `json:"data"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}

type SendRawTxArgs struct {
	Data string `json:"data"`
}

// type GetAccountArgs struct {
// 	Address string `json:"address"`
// }

type Transfer struct {
	txDB          badger.IStorage
	CurrentHeight string
	CurrentHash   string
}

type Txlog struct {
	Info          *SendTransaction
	TxHash        string
	CurrentHeight string
	CurrentHash   string
}

func NewTransfer(txDb badger.IStorage) *Transfer {
	return &Transfer{
		txDB: txDb,
	}
}

func sortAndEncodeMap(data map[string]string) string {
	mapkeys := make([]string, 0)
	for k := range data {
		mapkeys = append(mapkeys, k)
	}
	sort.Strings(mapkeys)
	strbuf := ""
	for i, key := range mapkeys {
		val := data[key]
		if val == "" {
			continue
		}
		strbuf += fmt.Sprintf("%s=%s", key, val)
		if i < len(mapkeys)-1 {
			strbuf += "&"
		}
	}
	return strbuf
}

func (t *Transfer) SignHash(tx *SendTransaction, key *ecdsa.PrivateKey) (string, error) {

	data := ""
	if tx.Data != "" && len(tx.Data) > 0 {
		data = "0x" + hex.EncodeToString([]byte(tx.Data))
	}

	tmp := map[string]string{
		"version":   config.Version,
		"to":        tx.To,
		"gas_price": tx.GasPrice,
		"gas_limit": tx.GasLimit,
		"data":      data,
		"nonce":     tx.Nonce,
		"value":     tx.Value,
	}
	enc := sortAndEncodeMap(tmp)

	if enc == "" {
		return "", fmt.Errorf("SignHash sort error")
	}
	hash := common.Bytes2Hash(ahash.SHA256([]byte(enc)))

	sig, err := crypto.ECDSASign(hash.Bytes(), key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sig), nil
}

// func (t *Transfer) SendTransactionFunc(args *SendRawTxArgs) (string, error) {
// 	var txhash string

// 	if err := t.updateWriteMsg(); err != nil {
// 		return "", err
// 	}

// 	if err := t.conn.CallMethod(1, "TxPool.SendRawTransaction", args, &txhash); err != nil {
// 		return "", err
// 	}
// 	return txhash, nil
// }

// func (t *Transfer) GetNonce(addr string) (string, error) {
// 	req := &GetAddrNonceByHashArgs{
// 		Address: addr,
// 	}
// 	var result string
// 	if err := t.conn.CallMethod(1, "TxPool.GetAddrTxNonce", &req, &result); err != nil {
// 		return "", err
// 	}
// 	return result, nil
// }

// func (t *Transfer) updateWriteMsg() error {
// 	head := make(map[string]interface{})
// 	err := t.conn.CallMethod(1, "Chain.Head", nil, &head)
// 	if err != nil {
// 		return err
// 	}

// 	t.CurrentHash = head["hash"].(string)
// 	t.CurrentHeight = fmt.Sprint(head["height"])
// 	return nil
// }

func (t *Transfer) WriteTxLog(txhash string, args *SendTransaction) error {

	txlog := &Txlog{
		TxHash:        txhash,
		CurrentHeight: t.CurrentHeight,
		CurrentHash:   t.CurrentHash,
		Info:          args,
	}

	bs, err := json.Marshal(txlog)
	if err != nil {
		return err
	}
	key := append(config.TxLogPrefix, []byte(txhash)...)
	return t.txDB.SetData(key, bs)
}

func (t *Transfer) ListTxLog() []string {
	result := []string{}
	t.txDB.PrefixForeachData(config.TxLogPrefix, func(k, v []byte) error {
		result = append(result, string(k))
		return nil
	})
	return result
}
