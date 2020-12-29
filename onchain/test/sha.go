/*
@Author : ZengLingHou
@Company :DingDu
@Description :
@Data: 2020/9/25 11:43
*/

package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"math/big"
	"unsafe"
)

func byte32(s []byte) (a *[32]byte) {
	if len(a) <= len(s) {
		a = (*[len(a)]byte)(unsafe.Pointer(&s[0]))
	}
	return a
}

func main() {
	h := sha3.NewKeccak256()
	temp, _ := new(big.Int).SetString("5604325519705830154105985269506265546337666327924488453635998662283027831438", 10)
	h.Write(abi.U256(temp))
	b := h.Sum(nil)
	hash := byte32(b)
	a := hash[:]
	result := hexutil.Encode(a)
	fmt.Println(result)
}
