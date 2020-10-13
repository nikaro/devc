package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// RemoveFromSlice return a slice of strings without the given string
func RemoveFromSlice(sliceIn []string, remove string) (sliceOut []string) {
	for _, item := range sliceIn {
		if item != remove {
			sliceOut = append(sliceOut, item)
		}
	}

	return sliceOut
}

// Md5Sum return md5 hash for a string
func Md5Sum(str string) (md5sum string) {
	hasher := md5.New()
	hasher.Write([]byte(str))
	md5sum = hex.EncodeToString(hasher.Sum(nil))

	return md5sum
}
