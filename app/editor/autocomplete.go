package editor

import (
	"fmt"
	"os"
	"strings"
)

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

func (ac *autoComplete) printWordsWithSharedPrefixes(prefix string) {
	fmt.Print("\n")

	words := ac.cmdTrie.getWordsGivenPrefix(prefix)
	for i, w := range words {
		fmt.Print(w)
		if i < len(words)-1 {
			fmt.Print("  ")
		}
	}
	fmt.Printf("\n")
}

func createCmdTrie() *Trie {
	trie := newTrie()

	builtinCommands := []string{"exit", "echo", "type", "pwd", "cd"}
	for _, name := range builtinCommands {
		trie.insert(name)
	}

	for dir := range strings.SplitSeq(os.Getenv("PATH"), ":") {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, file := range files {
			trie.insert(file.Name())
		}
	}

	return trie
}

type Trie struct {
	root *TrieNode
}

type TrieNode struct {
	children map[byte]*TrieNode
	isWord   bool
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		children: map[byte]*TrieNode{},
		isWord:   false,
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
	cur.isWord = true
}

func (t *Trie) complete(partialWord string) (string, int) {
	if len(partialWord) == 0 {
		return "", FOUND_NOTHING
	}

	cur := t.root

	for i := range partialWord {
		ch := partialWord[i]
		cur = cur.children[ch]
		if cur == nil {
			return "", FOUND_NOTHING
		}
	}

	var restOfWord strings.Builder
	for len(cur.children) != 0 {
		if cur.isWord || len(cur.children) > 1 {
			return restOfWord.String(), FOUND_MULTIPLE
		}
		// iterate once cause map has only one item
		for char, child := range cur.children {
			restOfWord.WriteByte(char)
			cur = child
			break
		}
	}

	return restOfWord.String(), FOUND_ONE
}

func (t *Trie) getWordsGivenPrefix(prefix string) []string {
	cur := t.root

	for i := range prefix {
		char := prefix[i]

		if child, found := cur.children[char]; found {
			cur = child
		} else {
			return nil // prefix doesn't exist in the trie
		}
	}

	result := []string{}
	builder := []byte(prefix)

	dfs(cur, &builder, &result)

	return result
}

func dfs(node *TrieNode, builder *[]byte, result *[]string) {
	if node.isWord {
		*result = append(*result, string(*builder))
	}

	for char, child := range node.children {

		*builder = append(*builder, char)

		dfs(child, builder, result)

		// backtrack
		*builder = (*builder)[:len(*builder)-1]
	}
}
