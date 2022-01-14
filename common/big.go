package common

import "math/big"

var (
	Big0          = new(big.Int).SetInt64(0)
	Big1          = new(big.Int).SetInt64(1)
	Big2          = new(big.Int).SetInt64(2)
	Big10         = new(big.Int).SetInt64(10)
	Big32         = new(big.Int).SetInt64(32)
	Big50         = new(big.Int).SetInt64(50)
	Big64         = new(big.Int).SetInt64(64)
	Big100        = new(big.Int).SetInt64(100)
	Big256        = new(big.Int).SetInt64(256)
	Big32Bits     = new(big.Int).Exp(Big2, Big32, nil)
	Big64Bits     = new(big.Int).Exp(Big2, Big64, nil)
	Big256Bits    = new(big.Int).Exp(Big2, Big256, nil)
	BigMaxUint32  = new(big.Int).Sub(Big32Bits, Big1)
	BigMaxUint64  = new(big.Int).Sub(Big64Bits, Big1)
	BigMaxUint256 = new(big.Int).Sub(Big256Bits, Big1)
	BigFloat0     = new(big.Float).SetInt(Big0)
	BigFloat1     = new(big.Float).SetInt(Big1)
)
