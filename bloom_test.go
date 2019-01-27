package bloom

import (
	"fmt"
	"testing"
)

func TestBloomHash(t *testing.T) {
	t.SkipNow()
	b := NewBloom(3, 3)
	for i := 0; i < 5; i++ {
		hh := b.hashes(i)
		fmt.Printf("\n%d:0\t%X", i, hh[0])
		fmt.Printf("\n%d:1\t%X", i, hh[1])
		fmt.Printf("\n%d:2\t%X", i, hh[2])
		fmt.Println()
	}
}

func TestBloomSetMember(t *testing.T) {
	BloomSetMember(NewBloom(20, 5), 10, t)
	BloomSetMember(NewBloom(200, 7), 500, t)
	BloomSetMember(NewBloom(200, 7), 10000, t)
	t.SkipNow()
	BloomSetMember(NewBloom(2000, 20), 100000, t)
}

func BloomSetMember(b Bloom, loops uint64, t *testing.T) {
	//t.SkipNow()
	col := uint64(0)
	fmt.Printf("\n\nBloom set with %v, %v - running for %v loops", b.n, b.k, loops)
	fmt.Printf("\nComputed FP proba\t: %01.5v", b.FalsePositiveProbability(loops))
	for i := uint64(0); i < loops; i++ {
		if b.Set(i) {
			col++
		}
	}

	fmt.Printf("\nActual collisions ratio\t: %01.5v", float64(col)/float64(loops))
	fmt.Printf("\nSimulated FP proba\t: %01.5v", b.FalsePositiveProbabilityEstimates())

	for i := uint64(0); i < loops; i++ {
		// Now, we expect all i to be members of the bloom filter ?
		if !b.Member(i) {
			//b.dump()
			t.Error(i) // BUG !
		}
	}
}

func TestTestBits(t *testing.T) {
	//t.SkipNow()
	b := NewBloom(20, 5)
	loops := 100
	for i := 0; i < loops; i++ {
		// Nothing set
		if b.testbit(uint64(i)) {
			fmt.Printf("\n%064b\n%064b\n", b.bf[0], b.bf[1])
			t.Error(i)
		}
		// Bits are set
		if b.setbit(uint64(i)) {
			fmt.Printf("\n%064b\n%064b\n", b.bf[0], b.bf[1])
			t.Error(i)
		}
		// Test bit are set
		if !b.testbit(uint64(i)) {
			fmt.Printf("\n%064b\n%064b\n", b.bf[0], b.bf[1])
			t.Error(i)
		}

		// Reset all bits
		if !b.setbit(uint64(i)) {
			fmt.Printf("\n%064b\n%064b\n", b.bf[0], b.bf[1])
			t.Error(i)
		}
	}
}

func TestTestBitsSparse(t *testing.T) {
	//t.SkipNow()
	b := NewBloom(2, 2)
	loops := 100
	for i := 0; i < loops; i++ {
		// Nothing set
		if b.testbit(uint64(i * 5)) {
			fmt.Printf("\n%064b\n%064b\n", b.bf[0], b.bf[1])
			t.Error(i)
		}
	}
	for i := 0; i < loops; i++ {
		// Bits are set
		if b.setbit(uint64(i * 5)) {
			fmt.Printf("\n%064b\n%064b\n", b.bf[0], b.bf[1])
			t.Error(i)
		}
	}
	for i := 0; i < loops; i++ {
		// Test bit are set
		if !b.testbit(uint64(i * 5)) {
			fmt.Printf("\n%064b\n%064b\n", b.bf[0], b.bf[1])
			t.Error(i)
		}
	}
	for i := 0; i < loops; i++ {
		// Re - set all bits
		if !b.setbit(uint64(i * 5)) {
			fmt.Printf("\n%064b\n%064b\n", b.bf[0], b.bf[1])
			t.Error(i)
		}
	}
}

func TestBloomSingleMember(t *testing.T) {

	b := NewBloom(3, 5)
	if b.Member(5) {
		t.Error()
	}
	if b.Set(5) {
		t.Error()
	}

	if !b.Member(5) {
		b.Dump()
		block, bit := hash2bits(b.hashes(5)[0], b.n)
		fmt.Printf("\nBlock=%v, bit=%v should have been set ?! ", block, bit)
		t.Error()
	} else {
		block, bit := hash2bits(b.hashes(5)[0], b.n)
		fmt.Printf("\nBlock=%v, bit=%v is set as expected ", block, bit)
	}
	if !b.Set(5) {
		b.Dump()
		block, bit := hash2bits(b.hashes(5)[0], b.n)
		fmt.Printf("\nBlock=%v, bit=%v should have been set ?! ", block, bit)
		t.Error()
	} else {
		block, bit := hash2bits(b.hashes(5)[0], b.n)
		fmt.Printf("\nBlock=%v, bit=%v is set as expected ", block, bit)
	}
	b.Dump()
}
