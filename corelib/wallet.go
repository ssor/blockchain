package corelib

import (
    "fmt"
    "crypto/ecdsa"
    "crypto/rand"
    "crypto/elliptic"
    "crypto/sha256"
    "golang.org/x/crypto/ripemd160"
)

func GetAllWallets() ([]*Wallet, error) {
    var wallets []*Wallet

    return wallets, nil
}

func AddWallet() (*Wallet, error) {
    return nil, fmt.Errorf("nil wallet")
}

func NewWallet() (*Wallet, error) {
    var w Wallet
    curve := elliptic.P256()
    private, err := ecdsa.GenerateKey(curve, rand.Reader)
    if err != nil {
        return nil, err
    } else {
        w.PrivateKey = private
    }
    w.PublicKey = append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
    return &w, nil
}

type Wallet struct {
    PrivateKey *ecdsa.PrivateKey
    PublicKey  []byte
}

func (wallet *Wallet) Address() ([]byte, error) {
    pubHash, err := hash160(wallet.PublicKey)
    if err != nil {
        return nil, err
    }
    versionPrePubHash:=append([]byte{byte(0x00)}, pubHash...)
    checksumData := checksum(versionPrePubHash,4)
    all:=append(versionPrePubHash,checksumData...)

}

func hash160(data []byte) ([]byte, error) {
    bs256 := sha256.Sum256(data)
    hash160 := ripemd160.New()
    _, err := hash160.Write(bs256[:])
    if err != nil {
        return nil, err
    }
    result := hash160.Sum(nil)
    return result, nil
}

func checksum(data []byte, length int) []byte {
    all := sha256.Sum256(sha256.Sum256(data)[:])
    if length > len(all) {
        return all[:]
    } else {
        return all[:length]
    }
}
