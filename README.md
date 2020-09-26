# bloom


[![PkgGoDev](https://pkg.go.dev/badge/github.com/oliverdaff/bloom)](https://pkg.go.dev/github.com/oliverdaff/bloom) [![Go Report Card](https://goreportcard.com/badge/github.com/oliverdaff/bloom)](https://goreportcard.com/report/github.com/oliverdaff/bloom) [![CircleCI](https://circleci.com/gh/oliverdaff/bloom.svg?style=shield)](https://circleci.com/gh/oliverdaff/bloom)

A simple bloom filter written in Go. 

A Bloom filter is a probabilistic data structure for set membership
that trades accuracy for space.  When queried for membership of a key a Bloom filter returns false if the key is definitely not in the set else it returns true which mean the key might be in the set.

This implementation of a Bloom filter uses a combination of the
murmur3 and fnv1 hashing to calculate which bits to set.

## API

__Create new Bloom filter__

A new BloomFilter is created using the `NewBloomFilter` function,
parameterized by:
*   `maxSize` - the maximum number of elements expected in the set.
*   `maxTolerance` - the expected accuracy (a sensible default is 0.01).
*   `seed` - the seed to use for the murmer hash function.

A error is returned if `maxSize` is set to 0 or the number of bits is needed in the backing bit set is larger than a uint32.
```go
bloom, err = NewBloomFilter(1000, 0.01, 42)
```

__Insert__

Insert a new key into the bloom filter using `Insert`.
```go
bloom.Insert([2,3,23,200])
```

__Contains__

Check if the key is contained in the set using `Contains`.

```go
bloom.Contains([23,89,205,148])
```



## Tests
The tests can be invoked with `go test`

## License
MIT © Oliver Daff