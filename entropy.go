package filetrove

import (
	"fmt"
	"io"
	"math"
	"os"
)

const (
	// MaxFileSize is the max size file that should be processed. This defaults to 1 GB.
	MaxFileSize = 1073741824
	// MaxEntropyChunk is the max byte size of a chunk read
	MaxEntropyChunk = 256000
)

// Entropy calculates the entropy of a file up to a hard-coded file size.
func Entropy(path string) (entropy float64, err error) {

	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	fStat, err := f.Stat()
	if err != nil {
		return 0, err
	}

	if !fStat.Mode().IsRegular() {
		return 0, fmt.Errorf("file (%s) is not a regular file to calculate entropy", path)
	}

	var filesize int64
	filesize = fStat.Size()
	if filesize == 0 {
		return 0, nil
	}

	if filesize > int64(MaxFileSize) {
		return 0, fmt.Errorf("file size (%d) is too large to calculate entropy (max allowed: %d)",
			filesize, int64(MaxFileSize))
	}

	dataBytes := make([]byte, MaxEntropyChunk)
	byteCounts := make([]int, 256)
	for {
		numBytesRead, err := f.Read(dataBytes)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		for i := 0; i < numBytesRead; i++ {
			byteCounts[int(dataBytes[i])]++
		}
	}

	for i := 0; i < 256; i++ {
		px := float64(byteCounts[i]) / float64(filesize)
		if px > 0 {
			entropy += -px * math.Log2(px)
		}
	}

	return math.Round(entropy*100) / 100, nil
}
