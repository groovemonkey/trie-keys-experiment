package main

import (
	"math/rand"
	"testing"

	"github.com/groovemonkey/trie-keys-experiment/mapkeys"
	"github.com/groovemonkey/trie-keys-experiment/prefix_trie_chunked"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ:."

// A helper for benchmarking this solution
func RandStringBytes(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

var realisticBenchmarkData = map[string]int{
	// Single metrics
	"foo":                                   1,
	"bar.baz":                               12,
	"business_summary.departments.finance":  13,
	"business_summary.departments.IT":       12,
	"business_summary.departments.software": 100,

	// Allows aggregating revenue
	"profits.revenue.top_line":               70,
	"profits.revenue.top_line.bake_sales":    70,
	"profits.revenue.top_line.charity":       70,
	"profits.revenue.top_line.asking_nicely": 70,

	"profits.revenue.top_line.enterprise_products":                                 0,
	"profits.revenue.top_line.enterprise_products.smalltime":                       70,
	"profits.revenue.top_line.enterprise_products.bigtime.but.not.that.big":        70,
	"profits.revenue.top_line.enterprise_products.this.is.a.long.subkey.wow.yeah":  17,
	"profits.revenue.top_line.enterprise_products.this.is.a.long.subkey.wow.no":    18,
	"profits.revenue.top_line.enterprise_products.this.is.a.long.subkey.wow.maybe": 19,

	"profits.revenue.top_line.cloud_tranformation": 70,
	"profits.revenue.taxes":                        -200,
	"profits.revenue.penalties":                    -50,
	"profits.revenue.fees":                         -10,

	"profits.revenue.bottom_line":                70,
	"profits.revenue.side_line":                  700,
	"profits.revenue.top_circle":                 403,
	"profits.revenue.let's.double_click_on_that": 7,
	"profits.revenue.net":                        3,
	"profits.revenue.basket":                     30,

	// very long and deeply nested prefixes
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_line":               70,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_line.bake_sales":    70,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_line.charity":       70,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_line.asking_nicely": 70,

	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_line.enterprise_products":                                 0,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_line.enterprise_products.smalltime":                       70,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_line.enterprise_products.bigtime.but.not.that.big":        70,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_line.enterprise_products.this.is.a.long.subkey.wow.yeah":  17,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_line.enterprise_products.this.is.a.long.subkey.wow.no":    18,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_line.enterprise_products.this.is.a.long.subkey.wow.maybe": 19,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_line.cloud_tranformation":                                 70,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.taxes":                                                        -200,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.penalties":                                                    -50,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.fees":                                                         -10,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.bottom_line":                                                  70,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.side_line":                                                    700,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.top_circle":                                                   403,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.let's.double_click_on_that":                                   7,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.net":                                                          3,
	"testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.sdf.revenue.basket":                                                       30,
}

// Make random data for test cases
func makeRandomDataMap(length int) map[int]string {
	// make a slice of test case strings, just long enough for all benchmark runs to complete
	testCases := make(map[int]string, length)

	for i := 0; i < length; i++ {
		randLen := randInt(1, 1000)
		randString := RandStringBytes(randLen)
		testCases[i] = randString
	}
	return testCases
}

// /////////////////
// // Maps
// /////////////////
func BenchmarkInsertMapRandom(b *testing.B) {
	data := makeRandomDataMap(b.N)
	store := make(mapkeys.Store[int])

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Insert(data[i], i)
	}
}

func BenchmarkInsertMapRealistic(b *testing.B) {
	store := make(mapkeys.Store[int])

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for key, val := range realisticBenchmarkData {
			store.Insert(key, val)
		}
	}
}

func BenchmarkSearchMapRandom(b *testing.B) {
	store := make(mapkeys.Store[int])
	// make a slice of test case strings, just long enough for all benchmark runs to complete
	data := makeRandomDataMap(b.N)

	// insert into store
	for val, key := range data {
		store.Insert(key, val)
	}

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// The function we're testing
		store.Search(data[i])
	}
}

func BenchmarkSearchMapRealistic(b *testing.B) {
	store := make(mapkeys.Store[int])

	// insert into store
	for key, val := range realisticBenchmarkData {
		store.Insert(key, val)
	}

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for key, _ := range realisticBenchmarkData {
			// The function we're testing
			store.Search(key)
		}
	}

}

func BenchmarkSearchPrefixMapRandom(b *testing.B) {
	store := make(mapkeys.Store[int])
	// make a slice of test case strings, just long enough for all benchmark runs to complete
	data := makeRandomDataMap(b.N)

	// Setup complete, let's bench
	var testString string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// chop the string in half
		testString = data[i]
		half := testString[:(len(testString) / 2)]
		testString = half
		b.StartTimer()

		// The function we're testing
		store.SearchPrefix(testString)
	}
}

func BenchmarkSearchPrefixMapRealistic(b *testing.B) {
	store := make(mapkeys.Store[int])

	// insert into store
	for key, val := range realisticBenchmarkData {
		store.Insert(key, val)
	}

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// short keys
		store.SearchPrefix("business_revenue")

		// medium keys
		store.SearchPrefix("profits")

		// long keys
		store.SearchPrefix("testing")

		// half of a medium key
		store.SearchPrefix("profits.revenue.top_line")

		// half of a wildly long key
		store.SearchPrefix("testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df")
	}
}

// /////////////////
// // Trie Chunked
// /////////////////
func BenchmarkInsertTrieChunkedRandom(b *testing.B) {
	store := prefix_trie_chunked.New[int]()

	// make a slice of test case strings, just long enough for all benchmark runs to complete
	data := makeRandomDataMap(b.N)

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Insert(data[i], i)
	}
}

func BenchmarkInsertTrieChunkedRealistic(b *testing.B) {
	store := prefix_trie_chunked.New[int]()

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for key, val := range realisticBenchmarkData {
			store.Insert(key, val)
		}
	}
}

func BenchmarkSearchTrieChunkedRandom(b *testing.B) {
	store := prefix_trie_chunked.New[int]()

	// make a slice of test case strings, just long enough for all benchmark runs to complete
	data := makeRandomDataMap(b.N)

	for val, key := range data {
		// insert into store
		store.Insert(key, val)
	}

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// The function we're testing
		store.Search(data[i])
	}
}

func BenchmarkSearchTrieChunkedRealistic(b *testing.B) {
	store := prefix_trie_chunked.New[int]()

	// insert into store
	for key, val := range realisticBenchmarkData {
		store.Insert(key, val)
	}

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for key, _ := range realisticBenchmarkData {
			// The function we're testing
			store.Search(key)
		}
	}

}

func BenchmarkSearchPrefixTrieChunkedRandom(b *testing.B) {
	store := prefix_trie_chunked.New[int]()

	// make a slice of test case strings, just long enough for all benchmark runs to complete
	data := makeRandomDataMap(b.N)

	for val, key := range data {
		// insert into store
		store.Insert(key, val)
	}

	// Setup complete, let's bench
	var testString string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// chop the string in half
		testString = data[i]
		half := testString[:(len(testString) / 2)]
		testString = half
		b.StartTimer()

		// The function we're testing
		store.SearchPrefix(testString)
	}
}

func BenchmarkSearchPrefixTrieChunkedRealistic(b *testing.B) {
	store := prefix_trie_chunked.New[int]()

	// insert into store
	for key, val := range realisticBenchmarkData {
		store.Insert(key, val)
	}

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// short keys
		store.SearchPrefix("business_revenue")

		// medium keys
		store.SearchPrefix("profits")

		// long keys
		store.SearchPrefix("testing")

		// half of a medium key
		store.SearchPrefix("profits.revenue.top_line")

		// half of a wildly long key
		store.SearchPrefix("testing.very.long.string.keys.with.many.many.many.many.segments.jhsdkfjhskdjhfks.kjhsdkjfhskdjhfksjhdf.kjshdkhskdjhfksjdhfkjshdfkjh.kjhsdkfjhskdjfhksjhdf.sd.sdf.sdf.sdf.sd.fs.dfs.dfs.dfs.df")
	}
}
