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

// HashSumsFile contains all hashes for a single file
type HashSumsFile struct {
	MD5        []byte
	SHA1       []byte
	SHA256     []byte
	SHA512     []byte
	BLAKE2B512 []byte
}

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

	checksum := hasher.Sum(nil)

	return checksum, nil
}

// HashAllFiles computes all supported hashes in a single file read using io.MultiWriter.
func HashAllFiles(inFile string) (HashSumsFile, error) {
	fd, err := os.Open(inFile)
	if err != nil {
		return HashSumsFile{}, err
	}
	defer fd.Close()

	md5h := md5.New()
	sha1h := sha1.New()
	sha256h := sha256.New()
	sha512h := sha512.New()
	blake2bh, err := blake2b.New512(nil)
	if err != nil {
		return HashSumsFile{}, err
	}

	mw := io.MultiWriter(md5h, sha1h, sha256h, sha512h, blake2bh)
	if _, err = io.Copy(mw, fd); err != nil {
		return HashSumsFile{}, err
	}

	return HashSumsFile{
		MD5:        md5h.Sum(nil),
		SHA1:       sha1h.Sum(nil),
		SHA256:     sha256h.Sum(nil),
		SHA512:     sha512h.Sum(nil),
		BLAKE2B512: blake2bh.Sum(nil),
	}, nil
}
