// Package hash contains Hash functions.
package hash

import (
	"context"
	"encoding/binary"
	"log"

	"github.com/minio/highwayhash"
	"github.com/mitchellh/hashstructure"
	"github.com/shengdoushi/base58"
)

var (
	highwayHashKey = []byte{83, 125, 180, 91, 99, 126, 30, 122, 153, 24, 56, 29, 78, 216, 80, 72, 214, 182, 101, 228, 170, 51, 229, 77, 58, 213, 68, 208, 68, 37, 154, 225}
)

// HighwayHash computes a fast non-cryptographic hash and returns a string.
func HighwayHash(ctx context.Context, obj interface{}) (hash string, err error) {
	highwayHash64, err := highwayhash.New64(highwayHashKey)
	if err != nil {
		log.Fatalln("Error initializing highway hash: " + err.Error())
	}
	highwayHashOptions := &hashstructure.HashOptions{
		Hasher:  highwayHash64,
		TagName: "hash",
	}

	hashUint64, err := hashstructure.Hash(obj, highwayHashOptions)
	if err != nil {
		return "", err
	}

	hashBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(hashBytes, hashUint64)

	return base58.Encode(hashBytes, base58.BitcoinAlphabet), nil
}

// HighwayHashUInt64 computes a fast non-cryptographic hash and returns it as an unsigned integer.
func HighwayHashUInt64(ctx context.Context, obj interface{}) (hash uint64, err error) {
	highwayHash64, err := highwayhash.New64(highwayHashKey)
	if err != nil {
		log.Fatalln("Error initializing highway hash: " + err.Error())
	}
	highwayHashOptions := &hashstructure.HashOptions{
		Hasher:  highwayHash64,
		TagName: "hash",
	}

	hashUint64, err := hashstructure.Hash(obj, highwayHashOptions)
	if err != nil {
		return 0, err
	}

	return hashUint64, nil
}
