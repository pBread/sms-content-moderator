package trie

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{children: make(map[rune]*TrieNode)}
}

// Insert a word into the trie
func (t *TrieNode) Insert(word string) {
	current := t
	for _, ch := range word {
		node, ok := current.children[ch]
		if !ok {
			node = NewTrieNode()
			current.children[ch] = node
		}
		current = node
	}
	current.isEnd = true // Mark the end of a word
}

// Search for a word in the trie
func (t *TrieNode) Search(word string) bool {
	current := t
	for _, ch := range word {
		node, ok := current.children[ch]
		if !ok {
			return false
		}
		current = node
	}
	return current.isEnd
}
