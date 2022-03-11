package dao

import (
	"crypto/ecdsa"
	"time"
	"walletboot/common"
	"walletboot/crypto"
	"walletboot/storage/badger"
)

var (
	addrKeyPre        = []byte("addr:")
	addrAccountPre    = []byte("addrobj:")
	addNewTime        = []byte("newtime:")
	defaultAddressKey = []byte("default")
)

type KeyStoreDB struct {
	storage *badger.Storage
}

func NewKeyStoreDB(storage *badger.Storage) *KeyStoreDB {
	return &KeyStoreDB{
		storage: storage,
	}
}

func (db *KeyStoreDB) GetDefaultAddress() (common.Address, error) {
	data, err := db.storage.GetData(defaultAddressKey)
	if err != nil {
		return common.Address{}, err
	}
	return common.Bytes2Address(data), nil
}

func (db *KeyStoreDB) Foreach(fn func(address common.Address, key *ecdsa.PrivateKey)) {
	_ = db.storage.PrefixForeachData(addrKeyPre, func(k []byte, v []byte) error {
		_, pkey, err := crypto.DecodePrivateKey(v)
		if err != nil {
			return err
		}
		addr := common.Bytes2Address(k)
		fn(addr, pkey)
		return nil
	})
}

func (db *KeyStoreDB) ForIndex(fn func(n int, k []byte, v []byte)) {
	db.ForIndex(fn)
}

func (db *KeyStoreDB) AddrForeach(fn func(k string, v []byte) error) error {
	return db.storage.PrefixForeach(string(addrAccountPre), fn)
}

func (db *KeyStoreDB) GetAddressNewTime(addr common.Address) ([]byte, error) {
	key := append(addNewTime, addr.Bytes()...)
	data, err := db.storage.GetData(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (db *KeyStoreDB) GetPrivateKey(address common.Address) (*ecdsa.PrivateKey, error) {
	key := append(addrKeyPre, address.Bytes()...)
	keyDer, err := db.storage.GetData(key)
	if err != nil {
		return nil, err
	}
	_, pkey, err := crypto.DecodePrivateKey(keyDer)
	if err != nil {
		return nil, err
	}
	return pkey, nil
}

func (db *KeyStoreDB) UpdateAccount(addr common.Address, obj []byte) error {
	sKey := append(addrAccountPre, addr.Bytes()...)
	return db.storage.SetData(sKey, obj)
}

func (db *KeyStoreDB) GetAccount(addr common.Address) ([]byte, error) {
	sKey := append(addrAccountPre, addr.Bytes()...)
	return db.storage.GetData(sKey)
}

func (db *KeyStoreDB) AccountsNumber() int64 {
	return db.storage.PrefixCount(string(addrAccountPre))
}

func (db *KeyStoreDB) Iterator() badger.Iterator {
	return db.storage.NewIteratorPrefix(addrAccountPre)
}

func (db *KeyStoreDB) Set(key, val []byte) error {
	return db.storage.SetData(key, val)
}

func (db *KeyStoreDB) Get(key []byte) ([]byte, error) {
	return db.storage.GetData(key)
}

func (db *KeyStoreDB) PutPrivateKey(addr common.Address, key *ecdsa.PrivateKey) error {
	sKey := append(addrKeyPre, addr.Bytes()...)
	keybytes := crypto.DefaultEncodePrivateKey(key)

	newTimeKey := append(addNewTime, addr.Bytes()...)
	newTime := time.Now().Unix()

	if err := db.storage.SetData(newTimeKey, common.Int2Byte(int(newTime))); err != nil {
		return err
	}
	if err := db.storage.SetData(sKey, keybytes); err != nil {
		return err
	}
	return nil
}

func (db *KeyStoreDB) SetDefaultAddress(address common.Address) error {
	return db.storage.SetData(defaultAddressKey, address.Bytes())
}

func (db *KeyStoreDB) RemoveAddress(address common.Address) error {
	key := append(addrKeyPre, address.Bytes()...)
	_, err := db.storage.GetData(key)
	if err != nil {
		return err
	}
	newTimeKey := append(addNewTime, address.Bytes()...)
	_, err = db.storage.GetData(newTimeKey)
	if err != nil {
		return err
	}

	if err := db.storage.DelData(key); err != nil {
		return err
	}
	if err := db.storage.DelData(newTimeKey); err != nil {
		return err
	}
	return nil
}

func (db *KeyStoreDB) DelDefault() error {
	addrByte, err := db.storage.GetData(defaultAddressKey)
	if err != nil {
		return err
	}

	addr := common.Bytes2Address(addrByte)

	newTimeKey := append(addNewTime, addr.Bytes()...)

	if err := db.storage.DelData(defaultAddressKey); err != nil {
		return err
	}
	if err := db.storage.DelData(newTimeKey); err != nil {
		return err
	}

	return nil
}
