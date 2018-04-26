package corelib

import (
    "fmt"
    "crypto/ecdsa"
    "crypto/rand"
    "crypto/elliptic"
    "crypto/sha256"
    "golang.org/x/crypto/ripemd160"
    "blockchain/db"
    "github.com/sirupsen/logrus"
    "os"
    pem2 "encoding/pem"
    "crypto/x509"
)

func GetAllWallets() ([]*Wallet, error) {
    var wallets []*Wallet
    db.GetDb().View(func(tx boltTx) error {
        b := tx.Bucket([]byte(db.BucketWallet))
        if b == nil {
            return nil
        }
        c := b.Cursor()
        for publicKey, prvKey := c.First(); publicKey != nil; publicKey, prvKey = c.Next() {
            wallet, err := decodeWallet(publicKey, prvKey)
            if err != nil {
                return err
            }
            wallets = append(wallets, wallet)
        }
        return nil
    })
    return wallets, nil
}

func AddWallet() (*Wallet, error) {
    w, err := NewWallet()
    if err != nil {
        logrus.Error("create wallet error: ", err)
        return nil, err
    }
    err = db.GetDb().Update(func(tx boltTx) error {
        bucket, err := tx.CreateBucketIfNotExists([]byte(db.BucketWallet))
        if err != nil {
            logrus.Errorf("create wallet bucket error: %s", err)
            return fmt.Errorf("cannot create db bucket")
        }
        prvKeyBs, err := w.privateKeySerialize()
        if err != nil {
            return err
        }
        //bs, err := w.serialize()
        //if err != nil {
        //    return err
        //}
        //address, err := w.Address()
        //if err != nil {
        //    return err
        //}
        err = bucket.Put(w.PublicKey, prvKeyBs)
        if err != nil {
            return err
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return w, fmt.Errorf("nil wallet")
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

func (wallet *Wallet) SaveToFile() (error) {
    address, err := wallet.Address()
    if err != nil {
        return err
    }
    privateKeyFileName := string(address) + "_private.pem"
    err = wallet.savePrivateKeyToFile(privateKeyFileName)
    if err != nil {
        return err
    }
    publicKeyFileName := string(address) + "_public.pem"
    err = wallet.savePublicKeyToFile(publicKeyFileName)
    if err != nil {
        return err
    }
    return nil
}

func (wallet *Wallet) privateKeySerialize() ([]byte, error) {
    return x509.MarshalECPrivateKey(wallet.PrivateKey)
}

func decodeWallet(pubKey, prvKey []byte) (*Wallet, error) {
    pKey, err := x509.ParseECPrivateKey(prvKey)
    if err != nil {
        return nil, err
    }
    return &Wallet{
        PublicKey:  pubKey,
        PrivateKey: pKey,
    }, nil
}

func (wallet *Wallet) savePrivateKeyToFile(fileName string) error {
    f, err := os.Create(fileName)
    if err != nil {
        return err
    }
    defer f.Close()

    privateBytes, err := x509.MarshalECPrivateKey(wallet.PrivateKey)
    if err != nil {
        return err
    }
    pemBlock := &pem2.Block{
        Type:  "Wallet Private",
        Bytes: privateBytes,
    }
    err = pem2.Encode(f, pemBlock)
    if err != nil {
        return err
    }
    return nil
}

func (wallet *Wallet) savePublicKeyToFile(fileName string) error {
    f, err := os.Create(fileName)
    if err != nil {
        return err
    }
    defer f.Close()

    pemBlock := &pem2.Block{
        Type:  "Wallet Private",
        Bytes: wallet.PublicKey,
    }
    err = pem2.Encode(f, pemBlock)
    if err != nil {
        return err
    }
    return nil
}

func (wallet *Wallet) Address() ([]byte, error) {
    pubHash, err := hash160(wallet.PublicKey)
    if err != nil {
        return nil, err
    }
    versionPrePubHash := append([]byte{byte(0x00)}, pubHash...)
    checksumData := checksum(versionPrePubHash, 4)
    all := append(versionPrePubHash, checksumData...)
    return base58Encode(all), nil
}

//func (wallet *Wallet) serialize() ([]byte, error) {
//    var result bytes.Buffer
//    encoder := gob.NewEncoder(&result)
//    err := encoder.Encode(wallet)
//    if err != nil {
//        blockLogger.Errorf("### encode error: %s\n", err)
//        return nil, fmt.Errorf("cannot serialize wallet")
//    }
//    return result.Bytes(), nil
//}
//
//func deserializeWallet(data []byte) (*Wallet, error) {
//    decoder := gob.NewDecoder(bytes.NewReader(data))
//    var wallet Wallet
//    err := decoder.Decode(&wallet)
//    if err != nil {
//        logrus.Errorf("### decode error: %s\n", err)
//        return nil, fmt.Errorf("cannot decode wallet")
//    } else {
//        return &wallet, nil
//    }
//}

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
    firstHash := sha256.Sum256(data)
    all := sha256.Sum256(firstHash[:])
    if length > len(all) {
        return all[:]
    } else {
        return all[:length]
    }
}
