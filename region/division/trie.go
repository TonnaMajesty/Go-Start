package main

type Division struct {
	codes []string
	info  Cell
}

type Cell struct {
	code string
	name string
}

type Node struct {
	children map[string]*Node
	value    *Division
}

func NewNode() *Node {
	return &Node{children: make(map[string]*Node)}
}

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	return &Trie{root: NewNode()}
}

func (t *Trie) Put(division *Division) {
	curr := t.root

	for _, r := range division.codes {
		if _, ok := curr.children[r]; !ok {
			curr.children[r] = NewNode()
		}
		curr = curr.children[r]
	}

	curr.value = division
}

func (t *Trie) Get(name []string) *Division {
	curr := t.root

	for _, r := range name {
		if _, ok := curr.children[r]; !ok {
			return nil
		}
		curr = curr.children[r]
	}

	return curr.value
}
