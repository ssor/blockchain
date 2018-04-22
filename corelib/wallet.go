package corelib

import "fmt"

func GetAllWallets() ([]*Wallet, error) {
    var wallets []*Wallet

    return wallets, nil
}

func AddWallet() (*Wallet, error) {

    return nil, fmt.Errorf("nil wallet")
}

type Wallet struct {

}
