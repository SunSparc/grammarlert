package main

// A multi-set implemented as a prefix tree (trie). Multi-sets allow
// the same key multiple times.

type Set struct {
	root  *setNode
	count int
}

type setNode struct {
	children map[byte]*setNode
	count    int
}

func NewSet() Set {
	return Set{
		root: newSetNode(),
	}
}

func newSetNode() *setNode {
	return &setNode{
		children: make(map[byte]*setNode),
	}
}

func (this *Set) Len() int {
	return this.count
}

func (this *Set) Put(key string) {
	node := this.root
	for i := 0; i < len(key); i++ {
		if _, exists := node.children[key[i]]; !exists {
			node.children[key[i]] = newSetNode()
		}
		node = node.children[key[i]]
	}
	node.count++
	this.count++
}

func (this *Set) Count(key string) int {
	node := this.getNode(key)
	if node == nil {
		return 0
	} else {
		return node.count
	}
}

func (this *Set) Has(key string) bool {
	node := this.getNode(key)
	if node == nil {
		return false
	} else {
		return node.count > 0
	}
}

func (this *Set) Find(prefix string, max int) []string {
	node := this.getNode(prefix)
	if node == nil {
		return []string{}
	}
	return this.findRecursively(prefix, max, node, []string{})
}

func (this *Set) findRecursively(prefix string, max int, node *setNode, results []string) []string {
	if node == nil {
		return results
	}

	if node.count > 0 {
		results = append(results, prefix)
	}

	if len(results) < max {
		for ch, node := range node.children {
			if node != nil {
				results = this.findRecursively(prefix+string(ch), max, node, results)
			}
		}
	}

	return results
}

func (this *Set) getNode(key string) *setNode {
	node := this.root
	for i := 0; i < len(key); i++ {
		if _, exists := node.children[key[i]]; !exists {
			return nil
		}
		node = node.children[key[i]]
	}
	return node
}
