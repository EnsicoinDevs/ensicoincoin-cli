package wallet

import (
	"crypto/sha256"
	"fmt"
	bolt "github.com/etcd-io/bbolt"
)

const (
	METADATA_BUCKET = "metadata"
)

type Wallet struct {
	path string
	db   *bolt.DB
}

func NewWallet(path string) *Wallet {
	return &Wallet{
		path: path,
	}
}

func (wallet *Wallet) Open() (err error) {
	wallet.db, err = bolt.Open(wallet.path, 0666, nil)
	if err != nil {
		return err
	}

	return wallet.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(METADATA_BUCKET))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		pkBytes := sha256.Sum256([]byte("olala"))

		return b.Put([]byte("pk"), pkBytes[:])
	})
}

func (wallet *Wallet) Close() error {
	return wallet.db.Close()
}
