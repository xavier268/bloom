# A Bloom filter implementation in go

## Usage

````
// Create a new bloom filter,
// 200 64bits words, using 20 hash functions
b := bloom.NewBloom(200, 20)

// Add some data to the filter
b.Set("Some data")
b.Set(65455)
b.Set(struct{765,"kjh"})

// Test for membership
if b.Member(testData) { ... }

````

## Foundation

You can find more details on the underlying theory on the excellent Wikipedia article on bloom filters.
