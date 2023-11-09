package mapkeys

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

func TestInsert(t *testing.T) {
	store := make(Store[int64])
	store.Insert("Hello, World!", 1)

	expected := "\nHello, World! - 1"
	result := store.String()
	if result != expected {
		t.Errorf("expected store.String() to be %s, but got %s", expected, result)
	}
}

func BenchmarkInsertMap(b *testing.B) {
	// make a slice of test case strings, just long enough for all benchmark runs to complete
	testCases := make(map[int]string, b.N)

	for i := 0; i < b.N; i++ {
		randLen := randInt(1, 1000)
		randString := RandStringBytes(randLen)
		testCases[i] = randString
	}

	store := make(Store[int])

	// Setup complete, let's bench
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.Insert(testCases[i], i)
	}
}

func TestSearch(t *testing.T) {
	store := make(Store[int64])

	store.Insert("aaa", 3)
	store.Insert("aa", 2)
	store.Insert("a", 1)

	found, val := store.Search("aaa")
	if (!found) || (val != 3) {
		t.Errorf("Expected to find val %d for a, found=%v, val=%d", 3, found, val)
	}
	found, val = store.Search("a")
	if (!found) || (val != 1) {
		t.Errorf("Expected to find val %d for a, found=%v, val=%d", 1, found, val)
	}
}

func BenchmarkSearchMap(b *testing.B) {
	store := make(Store[int])
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

func TestSearchPrefix(t *testing.T) {
	store := make(Store[int64])

	store.Insert("business_summary.departments.finance", 0)
	store.Insert("business_summary.departments.software", 100)
	store.Insert("business_summary.revenue.top_line", 70)
	store.Insert("business_summary.revenue.net", 50)

	results := store.SearchPrefix("business_summary")
	if len(results) == 0 {
		t.Errorf("Expected number of results to be > 0, store=%#v", store)
	}

	results = store.SearchPrefix("business_summary.departments")
	if len(results) != 2 {
		t.Errorf("Expected 2 results for query: results=%#v, store=%#v", results, store)
	}
}

func BenchmarkSearchPrefixMap(b *testing.B) {
	store := make(Store[int])
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
