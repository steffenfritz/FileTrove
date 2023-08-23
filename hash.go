package filetrove

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
	"io"
	"os"

	"golang.org/x/crypto/blake2b"
)

// ReturnSupportedHashes returns a list of supported hashes
func ReturnSupportedHashes() [5]string {
	return [5]string{"md5", "sha1", "sha256", "sha512", "blake2b-512"}
}

// Hashit hashes a file using the provided hash algorithm
func Hashit(inFile string, hashalg string) ([]byte, error) {
	fd, err := os.Open(inFile)
	if err != nil {
		return []byte{}, err
	}
	defer fd.Close()

	var hasher hash.Hash

	if hashalg == "sha256" {
		hasher = sha256.New()

	} else if hashalg == "md5" {
		hasher = md5.New()

	} else if hashalg == "sha1" {
		hasher = sha1.New()

	} else if hashalg == "sha512" {
		hasher = sha512.New()

	} else if hashalg == "blake2b-512" {
		hasher, err = blake2b.New512(nil)
		if err != nil {
			return []byte{}, err
		}

	} else {
		return []byte{}, errors.New("Hash not implemented.")
	}

	_, err = io.Copy(hasher, fd)
	if err != nil {
		return []byte{}, err
	}
	io.Copy(hasher, fd)

	checksum := hasher.Sum(nil)

	return checksum, nil
}
