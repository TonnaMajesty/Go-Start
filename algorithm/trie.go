package main

// Node represents a node in a trie
type Node struct {
	children map[rune]*Node // map of child nodes by rune
	value    interface{}    // optional value associated with the node
}

// NewNode creates a new node with an empty map of children
func NewNode() *Node {
	return &Node{children: make(map[rune]*Node)}
}

// Trie represents a trie data structure
type Trie struct {
	root *Node // root node of the trie
}

// NewTrie creates a new trie with an empty root node
func NewTrie() *Trie {
	return &Trie{root: NewNode()}
}

// Put inserts a key-value pair into the trie
func (t *Trie) Put(key string, value interface{}) {
	curr := t.root // start from the root node

	for _, r := range key { // iterate over each rune in the key
		if _, ok := curr.children[r]; !ok { // if the child node does not exist for this rune
			curr.children[r] = NewNode() // create a new child node for this rune
		}
		curr = curr.children[r] // move to the next child node
	}

	curr.value = value // assign the value to the last node
}

// Get returns the value associated with a key in the trie, or nil if not found
func (t *Trie) Get(key string) interface{} {
	curr := t.root // start from the root node

	for _, r := range key { // iterate over each rune in the key
		if _, ok := curr.children[r]; !ok { // if the child node does not exist for this rune
			return nil // return nil as the key is not found
		}
		curr = curr.children[r] // move to the next child node
	}

	return curr.value // return the value of the last node, or nil if no value is set
}
