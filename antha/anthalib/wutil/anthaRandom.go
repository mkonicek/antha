package wutil

import (
	"math/rand"
	"time"
)

// shared RNG

var ourRand *rand.Rand

// GetRandomWithSeed returns a random number generator, seeded with the seed provided if requested
func GetRandomWithSeed(seed int64, reseed bool) *rand.Rand {
	if reseed || ourRand == nil {
		ourRand = rand.New(rand.NewSource(seed))
	}
	return ourRand
}

// GetRandom returns a random number generator, seeding it if it is not initialised
func GetRandom() *rand.Rand {
	if ourRand == nil {
		ourRand = rand.New(rand.NewSource(time.Now().Unix()))
	}

	return ourRand
}
