// Package bloom implements a basic bloom filter
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
func NewBloom(n, k uint64) Bloom {
	return Bloom{n, k, make([]uint64, n)}
}

// Set a new member, in. Return true if already set
func (b *Bloom) Set(in interface{}) bool {
	r := true
	//fmt.Printf("\nSet %v", in)
	harr := b.hashes(in)
	//fmt.Printf("\nHash=%X", harr[0])
	for _, hh := range harr {
		r = r && b.setbit(hh)
		//fmt.Printf("\n%d\tbool=%v\t hash = %016X", i, r, hh)
	}
	return r
}

//Member returns true if the 'in' data is (mostlikely) already in the set.
func (b *Bloom) Member(in interface{}) bool {
	r := true
	//fmt.Printf("\nSet %v", in)
	harr := b.hashes(in)
	//fmt.Printf("\nHash=%X", harr[0])
	for _, hh := range harr {
		r = r && b.testbit(hh)
		//fmt.Printf("\n%d\tbool=%v\t hash = %016X", i, r, hh)
	}
	return r
}

// FalsePositiveCount give the expected max number of false positive for n members
// Estimates assumes the probability of false positive is low
func (b *Bloom) FalsePositiveCount(n uint64) uint64 {
	return uint64(b.FalsePositiveProbability(n) * float64(n))
}

// FalsePositiveProbability computes the probability of false positive for a given set
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

// ============= utilities =====================

//dump prints the internal structure of the filter
func (b *Bloom) dump() {
	fmt.Printf("\nNb of registers\t:%v", b.n)
	fmt.Printf("\nNb of hash func\t:%v", b.k)
	for i, r := range b.bf {
		fmt.Printf("\n%d\t%064b", i, r)
	}
	fmt.Println()
}

// Creates an array with the k hashed values
func (b *Bloom) hashes(in interface{}) []uint64 {
	h := fnv.New64()
	r := make([]uint64, b.k)
	for i := uint64(0); i < b.k; i++ {
		h.Reset()
		fmt.Fprintf(h, "==%d==%v==", b.k, in)
		r[i] = h.Sum64()
		//fmt.Printf("\n%X", r[i])
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
	//fmt.Printf("\nblock = %d, bit = %d", block, bit)
	var mask uint64 = 1 << bit
	//fmt.Printf("\nmask = \t%064b", mask)
	if b.bf[block]&mask == uint64(0) {
		b.bf[block] |= mask
		//fmt.Printf("\nbmap = \t%064b", b.bf[block])
		return false
	}
	return true
}
