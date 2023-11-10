package main

import (
	"fmt"

	"github.com/groovemonkey/trie-keys-experiment/mapkeys"
	"github.com/groovemonkey/trie-keys-experiment/prefix_trie"
)

func main() {
	mapStore := make(mapkeys.Store[int64])
	mapStore.Insert("one.two.three", 0)
	mapStore.Insert("one.two.three.one", 1)
	mapStore.Insert("one.two.three.two", 2)
	mapStore.Insert("one.two.three.two.one", 1)
	_, sum := mapStore.AggregateDescendants("one.two.three", mapkeys.Sum)

	fmt.Println("sum was:", sum)

	trie := prefix_trie.New[int64]()
	trie.Insert("Hello, world!", 42)
	trie.Insert("Hello, David!", 7)

	trie.Insert("aaaa", 4)
	trie.Insert("aa", 2)

	trie.Insert("business_summary.departments.finance", 0)
	trie.Insert("business_summary.departments.software", 100)
	trie.Insert("business_summary.revenue.top_line", 70)
	trie.Insert("business_summary.revenue.net", 50)

	found, node := trie.Search("floobtastic")
	if !found {
		fmt.Println("floobtastic wasn't in there")
		fmt.Println(node)
	}

	found, node = trie.Search("Hello, world!")
	if found {
		fmt.Printf("\nFound it! Val was %d\n", node.Value)
	}

	fmt.Print("\nTrying the prefix search with 'business_summary.revenue':")
	matches := trie.SearchPrefix("business_summary.revenue")
	for str, node := range matches {
		fmt.Printf("\n%s: %v", str, node.Value)
	}

	fmt.Println("\n\nPrinting trie")
	trie.DepthFirstPrint()
}
