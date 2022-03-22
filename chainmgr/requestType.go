package chainmgr

// type getBlockByHashArgs struct {
// 	Hash string `json:"hash"`
// }

type GetBlockHeaderByNumberArgs struct {
	Number string `json:"number"`
	//Count  string `json:"count"`
}

type getAccountArgs struct {
	RootHash string `json:"root_hash"`
	Address  string `json:"address"`
}

// type getAccountByAddrArgs struct {
// 	RootHash string `json:"root_hash"`
// 	Address  string `json:"address"`
// }

type SendRawTxArgs struct {
	Data string `json:"data"`
}

type GetAccountArgs struct {
	RootHash string `json:"root_hash"`
	Address  string `json:"address"`
}

// type GetAccountArgs struct {
// 	Address string `json:"address"`
// }

type GetTxByHashArgs struct {
	TxHash string `json:"hash"`
}

// type GetAccountArgs struct {
// 	RootHash string `json:"root_hash"`
// 	Address  string `json:"address"`
// }

type GetAddrNonceByHashArgs struct {
	Address string `json:"address"`
}
