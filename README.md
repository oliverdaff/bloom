# bloom

A simple bloom filter written in Go. 

A Bloom filter is a probabilistic data structure for set membership
that trades accuracy for space.  When queried for membership of a key 
a Bloom filter returns false if the key is definitely not in the set
else it returns true which mean the key might be in the set.



## Usage



## Tests
The tests can be invoked with `go test`

## License
MIT © Oliver Daff