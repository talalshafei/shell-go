package editor

const (
	FOUND_NOTHING = iota
	FOUND_ONE
	FOUND_MULTIPLE
)

type autoComplete struct {
	// all commands trie
	cmdTrie *Trie
}

func newAutoComplete() *autoComplete {
	cmdTrie := createCmdTrie()

	return &autoComplete{
		cmdTrie: cmdTrie,
	}
}

func (ac *autoComplete) completeWord(partialWord string) (string, int) {
	// first search inside path

	// search commands
	return ac.cmdTrie.complete(partialWord)
}

func createCmdTrie() *Trie {
	commandsNames := []string{"echo", "exit"}
	trie := newTrie()
	for _, name := range commandsNames {
		trie.insert(name)
	}

	return trie
}

type Trie struct {
	root *TrieNode
}

type TrieNode struct {
	children map[byte]*TrieNode
	word     string
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		children: map[byte]*TrieNode{},
		word:     "",
	}
}

func newTrie() *Trie {
	return &Trie{newTrieNode()}
}

func (t *Trie) insert(word string) {
	cur := t.root
	for i := range word {
		ch := word[i]
		if _, found := cur.children[ch]; !found {
			cur.children[ch] = newTrieNode()
		}
		cur = cur.children[ch]
	}
	cur.word = word
}

func (t *Trie) complete(partialWord string) (string, int) {
	cur := t.root

	for i := range partialWord {
		ch := partialWord[i]
		cur = cur.children[ch]
		if cur == nil {
			return "", FOUND_NOTHING
		}
	}

	var restOfWord []byte
	for len(cur.children) != 0 {
		if len(cur.children) > 1 {
			return "", FOUND_MULTIPLE
		}
		// iterate once cause map has only one item
		for char, next := range cur.children {
			restOfWord = append(restOfWord, char)
			cur = next
			break
		}
	}

	return string(restOfWord), FOUND_ONE
}
