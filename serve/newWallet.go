package serve

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"walletboot/common"
	"walletboot/config"
	"walletboot/httpxfs"
	"walletboot/storage/badger"

	"github.com/sirupsen/logrus"
)

type RandWallet struct {
	LoadAccountsDb badger.IStorage
	addrList       map[string]*Accounts
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
		addrList:       make(map[string]*Accounts),
	}
	return result
}

// 随机生成钱包地址
func (r *RandWallet) RandCreateWallet() error {

	addrList := make([]string, 0)
	if err := r.ClientHttp.CallMethod(1, "Wallet.List", nil, &addrList); err != nil {
		return err
	}
	fmt.Printf("addresslen:%v accountMax:%v\n", len(addrList), config.AccountNumberMax)
	if len(addrList) >= config.AccountNumberMax {
		return fmt.Errorf("AccountNumberMax:%v", config.AccountNumberMax)
	}
	var addr string
	if err := r.ClientHttp.CallMethod(1, "Wallet.Create", nil, &addr); err != nil {
		return err
	}

	r.addrList[addr] = &Accounts{
		Address: addr,
		Balance: "0",
	}
	return nil
}

// 随机获取钱包地址作为交易目标地址
func (r *RandWallet) GetTxTo() string {
	indexRand := rand.Intn(len(r.addrList))

	indexs := []string{}
	for index := range r.addrList {
		indexs = append(indexs, index)
	}

	for index, key := range indexs {
		if index == indexRand {
			return r.addrList[key].Address
		}
	}
	return ""
}

// 随机获取有额度的钱包作为交易的from
func (r *RandWallet) GetTxFrom() map[string]string {
	// update
	if err := r.TxFroms(); err != nil {
		logrus.Error(err)
		return nil
	}

	result := make(map[string]string)
	for addr, addrobj := range r.addrList {
		attobal, ok := new(big.Int).SetString(addrobj.Balance, 0)
		if !ok {
			addrobj.Balance = "0"
			continue
		}
		if attobal.Cmp(common.Big0) > 0 {
			result[addr] = addrobj.Balance
		}
	}
	return result
}

// 获取所有成为from资格的用户
func (r *RandWallet) TxFroms() error {

	addrList := make([]string, 0)
	if err := r.ClientHttp.CallMethod(1, "Wallet.List", nil, &addrList); err != nil {
		return err
	}

	var balance string

	for _, w := range addrList {
		if _, exits := r.addrList[w]; !exits {
			r.addrList[w] = &Accounts{
				Address: w,
				Balance: "0",
			}
		}
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
			r.addrList[w] = &Accounts{
				Address: w,
				Balance: balance,
			}
		}
	}
	return nil
}
