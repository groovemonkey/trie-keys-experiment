# Key/Value implementation Experimentation

Goal: an efficient representation of key/value pairs that supports the following operations:

- **Insert** - set a value for a given key
- **Search** - search for the value of a given key, if it exists
- **SearchPrefix** - search for all keys and values for which the key begins with this prefix

Different implementations are in different packages. `main` is just a playground.

## Testing
I'm writing basic tests as I go along, with most of the focus on benchmarking. To run tests and benchmarks:

```
go test -v ./... -bench=. -benchmem
```

