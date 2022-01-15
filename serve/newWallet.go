package serve

import (
	"encoding/json"
	"errors"
	"math/big"
	"math/rand"
	"testtx/common"
	"testtx/httpxfs"
	"testtx/storage/badger"

	"github.com/sirupsen/logrus"
)

type RandWallet struct {
	LoadAccountsDb badger.IStorage
	ClientHttp     *httpxfs.Client
}

type Accounts struct {
	Address string
	Balance string
}

func NewRqWallet(Db badger.IStorage, ClientHttp *httpxfs.Client) *RandWallet {
	result := &RandWallet{
		LoadAccountsDb: Db,
		ClientHttp:     ClientHttp,
	}
	return result
}

var AddrPrefix = []byte("addr:")
var AddrPrefixBal = []byte("addrbal:")

func (r *RandWallet) RandCreateWallet() error {
	var addr string
	if err := r.ClientHttp.CallMethod(1, "Wallet.Create", nil, &addr); err != nil {
		return err
	}

	key := string(AddrPrefix) + addr

	write := &Accounts{
		Address: addr,
		Balance: "0",
	}

	bs, _ := json.Marshal(write)
	return r.LoadAccountsDb.Set(key, bs)
}

func (r *RandWallet) GetTxTo() string {

	info := &Accounts{}
	star := r.LoadAccountsDb.PrefixCount(string(AddrPrefix))
	indexRand := rand.Intn(int(star))

	iter := r.LoadAccountsDb.NewIteratorPrefix(AddrPrefix)

	i := 0
	for iter.Next() {
		if i == indexRand {
			if err := json.Unmarshal(iter.Val(), &info); err != nil {
				logrus.Error(err)
				return ""
			}
		}
		i++
	}
	return info.Address
}

func (r *RandWallet) GetTxFrom() map[string]string {
	// update
	if err := r.TxFroms(); err != nil {
		logrus.Error(err)
		return nil
	}

	result := make(map[string]string)

	info := &Accounts{}
	star := r.LoadAccountsDb.PrefixCount(string(AddrPrefixBal))
	indexRand := rand.Intn(int(star))

	iter := r.LoadAccountsDb.NewIteratorPrefix(AddrPrefixBal)

	i := 0
	for iter.Next() {
		if i == indexRand {
			if err := json.Unmarshal(iter.Val(), &info); err != nil {
				logrus.Error(err)
				return nil
			}
		}
		i++
	}

	result[info.Address] = info.Balance
	return result
}

func (r *RandWallet) TxFroms() error {

	addrList := make([]string, 0)
	if err := r.ClientHttp.CallMethod(1, "Wallet.List", nil, &addrList); err != nil {
		return err
	}

	var balance string

	for _, w := range addrList {
		req := &GetAccountArgs{
			Address: w,
		}
		err := r.ClientHttp.CallMethod(1, "State.GetBalance", &req, &balance)
		if err != nil {
			return err
		}

		attobal, ok := new(big.Int).SetString(balance, 0)
		if !ok {
			return errors.New("func TxFroms string to big.Int err")
		}

		if attobal.Cmp(common.Big0) > 0 {

			key := string(AddrPrefixBal) + w
			info := &Accounts{
				Address: w,
				Balance: balance,
			}

			bs, _ := json.Marshal(info)
			if err := r.LoadAccountsDb.Set(key, bs); err != nil {
				return err
			}
		}
	}
	return nil
}
