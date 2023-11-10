package mapkeys

import (
	"testing"
)

func TestInsert(t *testing.T) {
	store := make(Store[int64])
	store.Insert("Hello, World!", 1)

	expected := "\nHello, World! - 1"
	result := store.String()
	if result != expected {
		t.Errorf("expected store.String() to be %s, but got %s", expected, result)
	}

	// duplicate should be fine
	store.Insert("dupecheck", 1)
	store.Insert("dupecheck", 1)
	found, val := store.Search("dupecheck")
	if (!found) || (val != 1) {
		t.Errorf("Expected duplicate Inserts to work. found=%v, val=%d", found, val)
	}

	// replacements should work
	store.Insert("replace", 100)
	store.Insert("replace", 300)
	found, val = store.Search("replace")
	if (!found) || (val != 300) {
		t.Errorf("Expected replacements to work. found=%v, val=%d", found, val)
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

func TestAggregateDescendants(t *testing.T) {
	mapStore := make(Store[int64])

	// test simple positives
	mapStore.Insert("one.two.three", 0)
	mapStore.Insert("one.two.three.one", 1)
	mapStore.Insert("one.two.three.two", 2)
	mapStore.Insert("one.two.three.two.one", 1)

	valid, sum := mapStore.AggregateDescendants("one.two.three", Sum)
	if !valid {
		t.Errorf("expected valid result")
	}
	if sum != 4 {
		t.Errorf("expected valid sum")
	}

	// test zero
	mapStore.Insert("floob", 0)
	valid, sum = mapStore.AggregateDescendants("floob", Sum)
	if !valid || (sum != 0) {
		t.Errorf("expected valid result of 0")
	}

	// test positive, negative, and zero numbers
	mapStore.Insert("negative", 0)
	mapStore.Insert("negative.one", -1)
	mapStore.Insert("negative.two", -2)
	mapStore.Insert("negative.two.zero", 0)
	mapStore.Insert("negative.two.one", 1)
	valid, sum = mapStore.AggregateDescendants("negative", Sum)
	if !valid || (sum != -2) {
		t.Errorf("expected valid result of -2, got %v", sum)
	}
}
