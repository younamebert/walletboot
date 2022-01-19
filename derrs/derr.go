package derrs

import "fmt"

var (
	CodeHttp    int64 = -10086
	CodeWallet  int64 = -20086
	CodeTranfer int64 = -30086
	// SystemErr   int64 = -40086
)

func NewErr(code int64, err string) error {
	return fmt.Errorf("%v:%v", code, err)
}
