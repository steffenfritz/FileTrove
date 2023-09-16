package filetrove

import (
	"bufio"
	"bytes"
	"errors"
	"go.etcd.io/bbolt"
	"os"
	"strings"
)

func CreateNSRLBoltDB(nsrlsourcefile string, nsrldbfile string) error {
	db, err := bbolt.Open(nsrldbfile, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	file, err := os.Open(nsrlsourcefile)
	if err != nil {
		return err
	}
	defer file.Close()

	batchSize := 100000
	values := make([]string, 0, batchSize)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hash := scanner.Text()
		values = append(values, hash)

		if len(values) == batchSize {
			err := db.Update(func(tx *bbolt.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists([]byte("sha1"))
				if err != nil {
					return err
				}

				for _, value := range values {
					err := bucket.Put([]byte(strings.ToLower(value)), []byte("true"))
					if err != nil {
						return err
					}
				}
				return nil
			})

			if err != nil {
				return err
			}
			values = values[:0]
		}
	}

	if len(values) > 0 {
		err := db.Update(func(tx *bbolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("sha1"))
			if err != nil {
				return err
			}

			for _, value := range values {
				err := bucket.Put([]byte(strings.ToLower(value)), []byte("true"))
				if err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// GetNSRL downloads the NSRL bolt database
func GetNSRLDB() {}

// ConnectNSRL connects to local bbolt NSRL file
func ConnectNSRL(nsrldbfile string) (*bbolt.DB, error) {
	db, err := bbolt.Open(nsrldbfile, 0600, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetValueNSRL reads bbolt database and checks if a given sha1 hash is present in the database
func GetValueNSRL(db *bbolt.DB, sha1hash []byte) (bool, error) {
	var fileIsInNSRL bool

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("sha1"))
		if b == nil {
			return errors.New("Could not connect to bucket.")
		}

		// the byte array translates to UTF-8 "true"
		fileIsInNSRL = bytes.Equal(b.Get(sha1hash), []byte{116, 114, 117, 101})
		// return nil to complete the transaction
		return nil
	})
	return fileIsInNSRL, err
}
