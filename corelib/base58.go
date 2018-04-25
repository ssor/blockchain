package corelib

import "math/big"

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func base58Encode(bs []byte) []byte {
    bi := big.NewInt(0).SetBytes(bs)
    bi0 := big.NewInt(0)
    base := big.NewInt(int64(len(b58Alphabet)))
    mod := &big.Int{}

    var result []byte
    for bi.Cmp(bi0) > 0 {
        bi.DivMod(bi, base, mod)
        result = append(result, b58Alphabet[mod.Int64()])
    }

    if bs[0] == 0x00 {
        result = append(result, b58Alphabet[0])
    }
    reverseBytes(result)
    return result
}

func reverseBytes(data []byte) {
    for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
        data[i], data[j] = data[j], data[i]
    }
}
