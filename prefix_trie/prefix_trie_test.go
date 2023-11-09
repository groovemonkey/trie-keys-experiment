package prefix_trie

import (
	"math/rand"
	"testing"
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

func BenchmarkInsertTrie(b *testing.B) {
	// make a slice of test case strings, just long enough for all benchmark runs to complete
	testCases := make(map[int]string, b.N)

	for i := 0; i < b.N; i++ {
		randLen := randInt(1, 1000)
		randString := RandStringBytes(randLen)
		testCases[i] = randString
	}

	store := New[int]()

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Insert(testCases[i], i)
	}
}

func BenchmarkSearchTrie(b *testing.B) {
	store := New[int]()
	// make a slice of test case strings, just long enough for all benchmark runs to complete
	testCases := make(map[int]string, b.N)

	// fill up the slice with test strings
	for i := 0; i < b.N; i++ {
		randLen := randInt(1, 1000)
		randString := RandStringBytes(randLen)
		// write to testcases so we can retrieve it later
		testCases[i] = randString

		// insert into store
		store.Insert(randString, i)
	}

	// Setup complete, let's bench
	var testString string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// every 100th search is for a (likely) nonexistent string
		if i%100 == 0 {
			randLen := randInt(1, 1000)
			testString = RandStringBytes(randLen)
		} else {

			testString = testCases[i]
		}
		b.StartTimer()

		// The function we're testing
		store.Search(testString)
	}
}

func BenchmarkSearchPrefixTrie(b *testing.B) {
	store := New[int]()
	// make a slice of test case strings, just long enough for all benchmark runs to complete
	testCases := make(map[int]string, b.N)

	// fill up the slice with test strings
	for i := 0; i < b.N; i++ {
		randLen := randInt(1, 1000)
		randString := RandStringBytes(randLen)
		// write to testcases so we can retrieve it later
		testCases[i] = randString

		// insert into store
		store.Insert(randString, i)
	}

	// Setup complete, let's bench
	var testString string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// every 100th search is for a (likely) nonexistent string
		if i%100 == 0 {
			randLen := randInt(1, 1000)
			testString = RandStringBytes(randLen)
		} else {
			// chop the string in half
			testString = testCases[i]
			half := testString[:(len(testString) / 2)]
			testString = half
		}
		b.StartTimer()

		// The function we're testing
		store.SearchPrefix(testString)
	}
}
