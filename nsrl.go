package filetrove

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	bbolt "go.etcd.io/bbolt"
	"os"
)

func CreateNSRLBoltDB(nsrlsourcefile string, nsrldbfile string) error {
	// Öffnen oder erstellen Sie die BoltDB-Datenbank
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

	batchSize := 50000
	values := make([]string, 0, batchSize)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hash := scanner.Text()
		values = append(values, hash)

		// Wenn die Batch-Größe erreicht ist, öffnen und schließen Sie die Transaktion
		if len(values) == batchSize {
			err := db.Update(func(tx *bbolt.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists([]byte("sha1"))
				if err != nil {
					return err
				}

				// Fügen Sie die Werte in die Datenbank ein
				for _, value := range values {
					err := bucket.Put([]byte(value), []byte("TRUE")) // Hier können Sie Ihre eigenen Daten speichern
					if err != nil {
						return err
					}
				}
				return nil
			})

			if err != nil {
				return err
			}
			values = values[:0] // Leeren Sie den Slice
		}
	}

	// Fügen Sie eventuell verbleibende Werte in die Datenbank ein
	if len(values) > 0 {
		err := db.Update(func(tx *bbolt.Tx) error {
			// Erstellen oder öffnen Sie den Eimer (Bucket)
			bucket, err := tx.CreateBucketIfNotExists([]byte("sha1"))
			if err != nil {
				return err
			}

			// Fügen Sie die verbleibenden Werte in die Datenbank ein
			for _, value := range values {
				err := bucket.Put([]byte(value), []byte("TRUE")) // Hier können Sie Ihre eigenen Daten speichern
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

	fmt.Println("Daten wurden erfolgreich hinzugefügt.")

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

		fileIsInNSRL = bytes.Equal(b.Get(sha1hash), []byte{116, 114, 117, 101})
		// return nil to complete the transaction
		return nil
	})
	return fileIsInNSRL, err
}
