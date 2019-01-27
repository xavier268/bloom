// Package bloom implements a bloom filter in pure go
package bloom

import (
	"fmt"
	"hash/fnv"
	"math"
)

// Bloom is the bloom filter object
type Bloom struct {
	n  uint64   // Nb of uint64 words
	k  uint64   // nb of hash functions
	bf []uint64 // internal data
}

// NewBloom provides a new, empty, bloom filter
// n is the number of 64bits words, k is the number of hash functions
func NewBloom(n, k uint64) Bloom {
	return Bloom{n, k, make([]uint64, n)}
}

// Set a new member. Return true if already set
func (b *Bloom) Set(member interface{}) bool {
	r := true
	//fmt.Printf("\nSet %v", in)
	harr := b.hashes(member)
	//fmt.Printf("\nHash=%X", harr[0])
	for _, hh := range harr {
		r = r && b.setbit(hh)
		//fmt.Printf("\n%d\tbool=%v\t hash = %016X", i, r, hh)
	}
	return r
}

//Member returns true if the member data is (most likely) already in the set.
// If it returns false, it is garanteed that member was not in the set.
func (b *Bloom) Member(member interface{}) bool {
	r := true
	harr := b.hashes(member)
	for _, hh := range harr {
		r = r && b.testbit(hh)
	}
	return r
}

// FalsePositiveProbability computes the probability of false positive for a given set size
// The formula does not seem to reflect the FalsePositiveProbabilityEstimates
// which is more precise, although longer to compute.
func (b *Bloom) FalsePositiveProbability(n uint64) (r float64) {
	r = 1 - 1/float64(b.n*64)
	r = 1 - math.Pow(r, float64(n*b.k))
	r = math.Pow(r, float64(b.k))
	return r
}

// FalsePositiveProbabilityEstimates simulates the FalsePositive proba on a sample
func (b *Bloom) FalsePositiveProbabilityEstimates(n uint64) float64 {
	bb := NewBloom(b.n, b.k)
	col := float64(0)
	for i := uint64(0); i < n; i++ {
		bb.Set(i)
	}
	for i := n; i < n+100000; i++ {
		if bb.Member(i) {
			col++
		}
	}
	return col / 100000.

}

//Dump prints the internal structure of the filter
func (b *Bloom) Dump() {
	fmt.Printf("\nNb of registers\t:%v", b.n)
	fmt.Printf("\nNb of hash func\t:%v", b.k)
	for i, r := range b.bf {
		fmt.Printf("\n%d\t%064b", i, r)
	}
	fmt.Println()
}

// ============= utilities =====================

// Creates an array with the k hashed values
func (b *Bloom) hashes(in interface{}) []uint64 {
	h := fnv.New64()
	r := make([]uint64, b.k)
	for i := uint64(0); i < b.k; i++ {
		h.Reset()
		fmt.Fprintf(h, "==%d==%v==", b.k, in)
		r[i] = h.Sum64()
	}
	return r
}

// given a hash, return the bit coordinate.
func hash2bits(h uint64, n uint64) (block uint64, bit uint64) {
	block = (h / uint64(64)) % n
	bit = h % uint64(64)
	return block, bit
}

// test if the selected bit is set, return true if set
func (b *Bloom) testbit(h uint64) bool {
	block, bit := hash2bits(h, b.n)
	var mask uint64 = 1 << bit
	if b.bf[block]&mask == mask {
		return true
	}
	return false
}

// Set the bit. Return true is already set.
func (b *Bloom) setbit(h uint64) bool {
	block, bit := hash2bits(h, b.n)
	var mask uint64 = 1 << bit
	if b.bf[block]&mask == uint64(0) {
		b.bf[block] |= mask
		return false
	}
	return true
}
