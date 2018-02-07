package wutil

import (
	"math/rand"
	"sync"
	"time"
)

// AnthaRandom is a thread-safe version of math/rand
type AnthaRandom struct {
	mutex sync.Mutex
	rng   *rand.Rand
}

// Intn returns, as an int, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
func (ar *AnthaRandom) Intn(n int) int {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Intn(n)
}

// Int returns a non-negative pseudo-random int.
func (ar *AnthaRandom) Int() int {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Int()
}

// Float64 returns, as a float64, a pseudo-random number in [0.0,1.0).
func (ar *AnthaRandom) Float64() float64 {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Float64()
}

// Float32 returns, as a float32, a pseudo-random number in [0.0,1.0).
func (ar *AnthaRandom) Float32() float32 {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Float32()
}

// ExpFloat64 returns an exponentially distributed float64 in the range (0, +math.MaxFloat64] with an exponential distribution whose rate parameter (lambda) is 1 and whose mean is 1/lambda (1).
func (ar *AnthaRandom) ExpFloat64() float64 {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.ExpFloat64()
}

// Int31 returns a non-negative pseudo-random 31-bit integer as an int32.
func (ar *AnthaRandom) Int31() int32 {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Int31()
}

// Int31n returns, as an int32, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
func (ar *AnthaRandom) Int31n(n int32) int32 {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Int31n(n)
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (ar *AnthaRandom) Int63() int64 {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Int63()
}

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n). It panics if n <= 0.
func (ar *AnthaRandom) Int63n(n int64) int64 {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Int63n(n)
}

// NormFloat64 returns a normally distributed float64 in the range [-math.MaxFloat64, +math.MaxFloat64] with standard normal distribution (mean = 0, stddev = 1).
func (ar *AnthaRandom) NormFloat64() float64 {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.NormFloat64()
}

// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers [0,n).
func (ar *AnthaRandom) Perm(n int) []int {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Perm(n)
}

// Read generates len(p) random bytes and writes them into p. It always returns len(p) and a nil error. Read should not be called concurrently with any other Rand method.
func (ar *AnthaRandom) Read(p []byte) (n int, err error) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Read(p)
}

// Seed uses the provided seed value to initialize the generator to a deterministic state. Seed should not be called concurrently with any other Rand method.
func (ar *AnthaRandom) Seed(seed int64) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	ar.rng.Seed(seed)
}

// Uint32 returns a pseudo-random 32-bit value as a uint32.
func (ar *AnthaRandom) Uint32() uint32 {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Uint32()
}

// Uint64 returns a pseudo-random 64-bit value as a uint64.
func (ar *AnthaRandom) Uint64() uint64 {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return ar.rng.Uint64()
}

// NewAnthaRandom returns a thread-safe RNG
func NewAnthaRandom(seed int64) *AnthaRandom {
	r := rand.New(rand.NewSource(seed))

	ar := AnthaRandom{rng: r}
	return &ar
}

// shared RNG
var ourRand *AnthaRandom

// GetRandomWithSeed returns a random number generator, seeded with the seed provided if requested
func GetRandomWithSeed(seed int64, reseed bool) *AnthaRandom {
	if reseed || ourRand == nil {
		ourRand = NewAnthaRandom(seed)
	}
	return ourRand
}

// GetRandom returns a random number generator, seeding it if it is not initialised
func GetRandom() *AnthaRandom {
	if ourRand == nil {
		ourRand = NewAnthaRandom(time.Now().Unix())
	}

	return ourRand
}
