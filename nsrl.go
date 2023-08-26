package filetrove

import bbolt "go.etcd.io/bbolt"

// GetNSRL downloads the NSRL sqlite database
func GetNSRLDB() {}

// ConnectNSRL connects to local bbolt NSRLR file
func ConnectNSRL(nsrldbfile string) (*bbolt.DB, error) {
	db, err := bbolt.Open(nsrldbfile, 0600, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetValueNSRL reads bbolt database and checks if a given sha1 hash is present in the database
func GetValueNSRL(sha1hash string) (bool, error) {}
