package wutil

import (
	"math/rand"
	"time"
)

// shared RNG

var ourRand *rand.Rand

func GetRandomWithSeed(seed int64) *rand.Rand {
	ourRand = rand.New(rand.NewSource(seed))
	return ourRand
}

func GetRandom() *rand.Rand {
	if ourRand == nil {
		ourRand = rand.New(rand.NewSource(time.Now().Unix()))
	}

	return ourRand
}
