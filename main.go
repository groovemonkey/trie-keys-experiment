package main

import "fmt"

type trieNode[T any] struct {
	Char  rune
	Value T
	// avoid mistaking initialized zero values for intentional zero values
	HasValue bool
	Children map[rune]*trieNode[T]
}

type Trie[T any] struct {
	root *trieNode[T]
}

func NewTrie[T any]() *Trie[T] {
	return &Trie[T]{root: &trieNode[T]{Children: make(map[rune]*trieNode[T])}}
}

func (t *Trie[T]) Insert(s string, val T) {
	currentNode := t.root
	for i, char := range s {
		child, ok := currentNode.Children[char]
		// If there's no such child, create one
		if !ok {
			child = &trieNode[T]{Char: char, Children: make(map[rune]*trieNode[T])}
			currentNode.Children[char] = child
		}
		// is this a leaf node? If so, set the value.
		if i == (len(s) - 1) {
			child.Value = val
			child.HasValue = true
		}
		// set currentNode for the next iteration
		currentNode = child
	}
}

func (t *Trie[T]) DepthFirstPrint() {
	if t.root == nil {
		return
	}
	currentNode := t.root
	depthFirstPrint(currentNode, "")
}

// Search returns whether or not the search string exists in the Trie, and if it does, the associated node.
func (t *Trie[T]) Search(s string) (bool, *trieNode[T]) {
	var returnVal *trieNode[T]

	currentNode := t.root
	for _, char := range s {
		child, ok := currentNode.Children[char]
		if !ok {
			return false, returnVal
		}
		currentNode = child
	}
	// we're at the last character, it's a match
	// NOTE: we don't care whether this is a valid key (i.e. whether currentNode.HasValue)
	return true, currentNode
}

// SearchPrefix returns all (string) Keys with the given prefix, mapped to pointers to their value-containing trieNodes
// a "Key" here means a trieNode that has a value associated with it, i.e. the last character of a key
func (t *Trie[T]) SearchPrefix(prefix string) map[string]*trieNode[T] {
	keysAndVals := make(map[string]*trieNode[T])
	if len(prefix) == 0 {
		return keysAndVals
	}
	found, node := t.Search(prefix)
	if !found {
		fmt.Println("SearchPrefix: nothing found.")
		return keysAndVals
	}
	// find all descendants of the node, after trimming the prefix
	// NOTE: we trim the prefix because getDescendants() would duplicate the first letter of the matched prefix
	// (it immediately adds the current Node's rune to the prefix, which would duplicate the last rune)
	trimmedPrefix := string([]rune(prefix)[:len(prefix)-1])
	return getDescendants[T](node, trimmedPrefix, make(map[string]*trieNode[T]))
}

// getDescendants is a depth-first search starting at a node and returning a slice of descendant Nodes that represent a valid Key (they have a Value)
// TODO(dcohen) make this faster (benchmark!) by passing the matchedNodes map by pointer instead of by value
func getDescendants[T any](currentNode *trieNode[T], prefix string, matchedNodes map[string]*trieNode[T]) map[string]*trieNode[T] {
	stringUntilNow := fmt.Sprintf("%s%c", prefix, currentNode.Char)

	// Are we a node that contains a value? (end of a Key?)
	if currentNode.HasValue {
		matchedNodes[stringUntilNow] = currentNode
	}

	for _, node := range currentNode.Children {
		for key, node := range getDescendants[T](node, stringUntilNow, matchedNodes) {
			matchedNodes[key] = node
		}
	}
	return matchedNodes
}

// This works.
// func depthFirst[T any](currentNode *trieNode[T]) {
// 	fmt.Printf("%c", currentNode.Char)
//
// 	// obvious problem with intentional zero values for ints
// 	if currentNode.HasValue {
// 		fmt.Printf("\nValue: %v \n", currentNode.Value)
// 	}
// 	for _, node := range currentNode.Children {
// 		depthFirst[T](node)
// 	}
// }

// depthFirstPrint with accumulator
// TODO(dcohen) return a string here, by doing the normal recursive "return acc + depthFirstPrint(...)"
func depthFirstPrint[T any](currentNode *trieNode[T], acc string) {
	stringUntilNow := fmt.Sprintf("%s%c", acc, currentNode.Char)

	// Is this the last rune of a Key?
	if currentNode.HasValue {
		fmt.Printf("Key: %s Value: %v\n", stringUntilNow, currentNode.Value)
	}
	for _, node := range currentNode.Children {
		depthFirstPrint[T](node, stringUntilNow)
	}
}

func main() {
	trie := NewTrie[int64]()
	trie.Insert("Hello, world!", 42)
	trie.Insert("Hello, David!", 7)

	// trie.Insert("b", 1)
	// trie.Insert("bb", 2)
	// trie.Insert("bbb", 3)

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
		fmt.Printf("\nFound it! Val was %d \n", node.Value)
	}

	fmt.Println("Trying the prefix search with 'Hello'")
	matches := trie.SearchPrefix("business_summary.revenue")
	for str, node := range matches {
		fmt.Printf("\n%s: %v", str, node.Value)
	}

	fmt.Println("\n\nPrinting trie")
	trie.DepthFirstPrint()
}
