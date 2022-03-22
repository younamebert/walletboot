package serve

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sync"
	"time"
	"walletboot/common"
	"walletboot/config"
	"walletboot/crypto"
	"walletboot/dao"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// Wallet represents a software wallet that has a default address derived from private key.
type Wallet struct {
	db          *dao.KeyStoreDB
	mu          sync.RWMutex
	cacheMu     sync.RWMutex
	defaultAddr common.Address
	cache       map[common.Address]*ecdsa.PrivateKey
}

type Accounts struct {
	Address   string
	Balance   decimal.Decimal `json:"balance"`
	Nonce     uint64          `json:"nonce"`
	Extra     string          `json:"extra"`
	Code      string          `json:"code"`
	StateRoot string          `json:"state_root"`
}

// NewWallet constructs and returns a new Wallet instance with badger db.
func NewWallet(daokeyDB *dao.KeyStoreDB) (*Wallet, error) {

	w := &Wallet{
		db:    daokeyDB,
		cache: make(map[common.Address]*ecdsa.PrivateKey),
		// conn:  cli,
	}

	if _, err := w.db.GetDefaultAddress(); err != nil {
		coinbase, err := w.AddByRandom()
		if err != nil {
			return nil, err
		}
		w.defaultAddr = coinbase
		if err = w.db.SetDefaultAddress(coinbase); err != nil {
			return nil, err
		}
	}

	if err := w.setupTxFrom(); err != nil {
		return nil, err
	}
	return w, nil
}

//There must be an initialized from
func (w *Wallet) setupTxFrom() error {
	addrdef, err := w.db.GetDefaultAddress()
	if err != nil {
		return err
	}
	w.defaultAddr = addrdef
	defaultAddrAccount, err := w.GetAccount(addrdef)
	if err != nil {
		createWallet := &Accounts{
			Address: addrdef.B58String(),
			Balance: decimal.NewFromInt(0),
		}
		// json.Marshal(createWallet)
		bs, _ := json.Marshal(createWallet)
		w.db.UpdateAccount(addrdef, bs)
		return fmt.Errorf("2 no user with from qualification was found in the wallet launcher. Please add XFS coins to this address:%v", addrdef.String())
	}
	if defaultAddrAccount.Balance == decimal.NewFromInt(0) {
		return fmt.Errorf("3 no user with from qualification was found in the wallet launcher. Please add XFS coins to this address:%v", addrdef.B58String())
	}
	return nil
}

// Blocks that were successfully traded in the previous block
// func (w *Wallet) UpdateAccount() error {

// }

func (w *Wallet) Accounts() {
	iter := w.db.Iterator()

	accounts := make([]*Accounts, 0)
	ait := &Accounts{}
	for iter.Next() {
		if err := json.Unmarshal(iter.Val(), ait); err != nil {
			continue
		}
		accounts = append(accounts, ait)
	}
	bs, _ := common.MarshalIndent(accounts)
	fmt.Println(string(bs))
}

// AddByRandom constructs a new Wallet with a random number and retuens the its address.
func (w *Wallet) AddByRandom() (common.Address, error) {
	key, err := crypto.GenPrvKey()
	if err != nil {
		return common.Address{}, err
	}
	return w.AddWallet(key)
}

func (w *Wallet) GetWalletNewTime(addr common.Address) ([]byte, error) {
	return w.db.GetAddressNewTime(addr)
}

func (w *Wallet) CreateAccount() (common.Address, error) {
	// max account number
	if w.GetNumber() >= config.AccountMaxNumber {
		return common.Address{}, fmt.Errorf("create account max number:%v", config.AccountMaxNumber)
	}

	addr, err := w.AddByRandom()
	if err != nil {
		return common.Address{}, err
	}

	accounts := &Accounts{
		Balance: decimal.NewFromInt(0),
		Address: addr.B58String(),
	}

	bs, err := json.Marshal(accounts)
	if err != nil {
		return common.Address{}, err
	}
	if err := w.db.UpdateAccount(addr, bs); err != nil {
		return common.Address{}, err
	}
	return addr, nil
}

//随机获取账户
func (w *Wallet) RandomAccessAccount() []*Accounts {
	var (
		accounts []*Accounts = make([]*Accounts, 0)
	)
	accountLength := w.db.AccountsNumber()
	if accountLength == 0 {
		return accounts
	}
	rand.Seed(time.Now().UnixNano())
	accountIndex := rand.Intn(int(accountLength) - 1)
	froms := make([]*Accounts, 0)

	i := 0
	err := w.db.AddrForeach(func(k string, v []byte) error {
		account := &Accounts{}
		if err := json.Unmarshal(v, &account); err != nil {
			return err
		}
		if i == accountIndex {
			if err := json.Unmarshal(v, &account); err != nil {
				return err
			}
			accounts = append(accounts, account)
		} else {
			if account.Balance.BigInt().Cmp(big.NewInt(0)) > 0 {
				froms = append(froms, account)
			}
		}
		i++
		return nil
	})

	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}
	if len(froms) == 0 {
		return froms
	}
	accountIndexFrom := rand.Intn(len(froms))
	accounts = append(accounts, froms[accountIndexFrom])
	return accounts
}

func (w *Wallet) GetNumber() int {
	i := 0
	iter := w.db.Iterator()
	for iter.Next() {
		i++
	}
	return i
}

func (w *Wallet) GetAccount(addr common.Address) (*Accounts, error) {
	val, err := w.db.GetAccount(addr)
	if err != nil {
		return nil, err
	}

	accounts := &Accounts{}
	if err := json.Unmarshal(val, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (w *Wallet) AddWallet(key *ecdsa.PrivateKey) (common.Address, error) {
	addr := crypto.DefaultPubKey2Addr(key.PublicKey)
	if err := w.db.PutPrivateKey(addr, key); err != nil {
		return common.Address{}, err
	}
	if w.defaultAddr.Equals(common.Address{}) {
		if err := w.SetDefault(addr); err != nil {
			return addr, nil
		}
	}
	return addr, nil
}

func (w *Wallet) All() map[common.Address]*ecdsa.PrivateKey {
	data := make(map[common.Address]*ecdsa.PrivateKey)
	w.db.Foreach(func(address common.Address, key *ecdsa.PrivateKey) {
		data[address] = key
	})
	return data
}

func (w *Wallet) GetKeyByAddress(address common.Address) (*ecdsa.PrivateKey, error) {
	w.cacheMu.RLock()
	if pk, has := w.cache[address]; has {
		w.cacheMu.RUnlock()
		return pk, nil
	}
	w.cacheMu.RUnlock()
	key, err := w.db.GetPrivateKey(address)
	if err != nil {
		return nil, err
	}
	w.cacheMu.Lock()
	w.cache[address] = key
	w.cacheMu.Unlock()
	return key, nil
}

func (w *Wallet) SetDefault(address common.Address) error {
	if address.Equals(w.defaultAddr) {
		return nil
	}
	k, err := w.GetKeyByAddress(address)
	if err != nil || k == nil {
		return fmt.Errorf("not found address %s", address.B58String())
	}
	err = w.db.SetDefaultAddress(address)
	if err != nil {
		return err
	}
	w.mu.Lock()
	w.defaultAddr = address
	w.mu.Unlock()
	return nil
}

func (w *Wallet) GetDefault() common.Address {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.defaultAddr
}

func (w *Wallet) Remove(address common.Address) error {

	if address.Equals(w.defaultAddr) {
		return errors.New("default address cannot be deleted")
		// w.mu.Lock()
		// if err := w.db.DelDefault(); err != nil {
		// 	w.mu.Unlock()
		// 	return err
		// }
		// w.defaultAddr = noneAddress

		// w.mu.Unlock()
	}
	w.mu.Lock()
	if err := w.db.RemoveAddress(address); err != nil {
		return err
	}
	w.mu.Unlock()
	w.cacheMu.Lock()
	delete(w.cache, address)
	w.cacheMu.Unlock()
	return nil
}

func (w *Wallet) Export(address common.Address) ([]byte, error) {
	key, err := w.GetKeyByAddress(address)
	if err != nil {
		return nil, err
	}
	return crypto.DefaultEncodePrivateKey(key), nil
}

func (w *Wallet) Import(der []byte) (common.Address, error) {
	kv, pKey, err := crypto.DecodePrivateKey(der)
	if err != nil {
		return common.Address{}, err
	}
	if kv != crypto.DefaultKeyPackVersion {
		return common.Address{}, fmt.Errorf("unknown private key version %d", kv)
	}
	return w.AddWallet(pKey)
}
