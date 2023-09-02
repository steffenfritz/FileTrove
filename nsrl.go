package filetrove

import (
	"bufio"
	"bytes"
	"errors"
	bbolt "go.etcd.io/bbolt"
	"os"
)

// CreateNSRLBoltDB creates a BoltDB file from a specific text file that contains one hash sum per line
func CreateNSRLBoltDB(nsrlsourcefile string, nsrldbfile string) error {

	// Create db file
	db, err := bbolt.Open(nsrldbfile, 0600, &bbolt.Options{})
	if err != nil {
		return err
	}
	defer db.Close()

	// Create bucket
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("sha1"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Read the NSRL source file where each line MUST contain a single SHA1 hash
	readFile, err := os.Open(nsrlsourcefile)
	if err != nil {
		return err
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		// Add hashes to the sha1 bucket
		db.Update(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte("sha1"))
			if bucket == nil {
				return errors.New("Could not connect to bucket.")
			}
			return bucket.Put([]byte(fileScanner.Text()), []byte("true"))
		})
	}

	err = readFile.Close()

	return err
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

		fileIsInNSRL = bytes.Equal(b.Get(sha1hash), []byte{116, 114, 117, 101})
		// return nil to complete the transaction
		return nil
	})
	return fileIsInNSRL, err
}
