# Key/Value implementation Experimentation

Goal: an efficient representation of key/value pairs that supports the following operations:

- **Insert** - set a value for a given key
- **Search** - search for the value of a given key, if it exists
- **SearchPrefix** - search for all keys and values for which the key begins with this prefix

Different implementations are in different packages. `main` is just a playground.

## Implementations

1. **Map** - A simple map implementation.
1. **Prefix Trie** -- A simple trie that stores one character per node. E.g. The key "hello.world" is stored as 11 nodes.
1. **Chunked Prefix Trie** -- A slightly less abstract implementation that stores a dot-separated key chunk (a string, not just one character/codepoint) in each node to minimize pointer chasing for longer strings. E.g. "hello.world" is stored as 2 nodes.

## Testing
I'm writing basic tests as I go along, with most of the focus on benchmarking. To run tests and benchmarks:

```
go test -v ./... -bench=. -benchmem
```

## Benchmarks (updated Nov 9, 2023)

Conclusions:
- maps are blazing fast at insertion and search, as expected.
- maps are slow for SearchPrefix
- surprisingly, the simple trie is insanely slow, probably due to excessive pointer chasing (and longer strings, on average)
- the chunked trie is really good:
    - 4x slower than the map on Insert,
    - 3x slower on Search,
    - 45x faster on SearchPrefix

### Map Implementation

```
pkg: github.com/groovemonkey/trie-keys-experiment/mapkeys
BenchmarkInsertMap
BenchmarkInsertMap-12          	 2513421	       449.9 ns/op	      96 B/op	       0 allocs/op
BenchmarkSearchMap
BenchmarkSearchMap-12          	 3007292	       424.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchPrefixMap
BenchmarkSearchPrefixMap-12    	   10000	    121197 ns/op	     902 B/op	       0 allocs/op
```

### Trie (one rune per node)

```
pkg: github.com/groovemonkey/trie-keys-experiment/prefix_trie
BenchmarkInsertTrie
BenchmarkInsertTrie-12          	   26985	     43263 ns/op	   95877 B/op	    1498 allocs/op
BenchmarkSearchTrie
BenchmarkSearchTrie-12          	  189331	     31631 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchPrefixTrie
BenchmarkSearchPrefixTrie-12    	   10000	    391979 ns/op	  186139 B/op	     768 allocs/op
```

### Trie (chunked)

```
pkg: github.com/groovemonkey/trie-keys-experiment/prefix_trie_chunked
BenchmarkInsertTrieChunked
BenchmarkInsertTrieChunked-12          	  569263	      1852 ns/op	    3169 B/op	      30 allocs/op
BenchmarkSearchTrieChunked
BenchmarkSearchTrieChunked-12          	  845596	      1299 ns/op	     165 B/op	       1 allocs/op
BenchmarkSearchPrefixTrieChunked
BenchmarkSearchPrefixTrieChunked-12    	  715863	      2711 ns/op	    2245 B/op	      17 allocs/op
```

All benchmarks:
```
goos: darwin
goarch: arm64
```

