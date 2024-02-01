package hash

import (
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

func Keccak256(input []byte) common.Hash {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(input)
	return common.BytesToHash(hash.Sum(nil))
}
