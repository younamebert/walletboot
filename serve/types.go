package serve

type getAccountArgs struct {
	RootHash string `json:"root_hash"`
	Address  string `json:"address"`
}
