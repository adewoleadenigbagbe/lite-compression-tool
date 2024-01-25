package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

type node struct {
	char      byte
	frequency int64
}

type nodeArray []node

func (n nodeArray) Len() int           { return len(n) }
func (n nodeArray) Less(i, j int) bool { return n[i].frequency < n[j].frequency }
func (n nodeArray) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n nodeArray) Sort()              { sort.Sort(n) }

type HuffmanNode struct {
	left   *HuffmanNode
	right  *HuffmanNode
	char   byte
	weight int64
}

func getFile(filename string) (*os.File, error) {
	return os.Open(filename)
}

func calculateCharacterFrequency(file *os.File) (map[byte]int64, error) {
	const chunkSize = 10000
	buffer := make([]byte, chunkSize)
	frequencyTable := make(map[byte]int64)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err == io.EOF {
			break
		}

		for _, v := range buffer[:n] {
			frequencyTable[v]++
		}
	}

	return frequencyTable, nil
}

func BuildHuffmanTree(nodes nodeArray) *HuffmanNode {
	root := &HuffmanNode{}

	for _, n := range nodes {
		h := &HuffmanNode{
			char:   n.char,
			weight: n.frequency,
		}

		if root.left == nil {
			root.left = h
			root.weight += h.weight
			continue
		} else if root.right == nil {
			root.right = h
			root.weight += h.weight
			continue
		}

		current := &HuffmanNode{
			weight: root.weight + h.weight,
		}

		if root.weight < h.weight {
			current.left = root
			current.right = h
		} else {
			current.right = root
			current.left = h
		}

		root = current
	}

	return root
}

func PreOrderTraversal(root *HuffmanNode) {
	if root == nil {
		return
	}
	fmt.Print(root.weight, " ")
	PreOrderTraversal(root.left)
	PreOrderTraversal(root.right)
}

func main() {
	var (
		err            error
		file           *os.File
		frequencyTable = make(map[byte]int64)
	)

	file, err = getFile("lesmiserables.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	frequencyTable, err = calculateCharacterFrequency(file)
	if err != nil {
		log.Fatal(err)
	}

	var nodes nodeArray
	for k, v := range frequencyTable {
		node := node{
			char:      k,
			frequency: v,
		}
		nodes = append(nodes, node)
	}

	if len(nodes) < 2 {
		fmt.Println("Huffman tree cannot be built with less than two nodes")
		return
	}

	//sort the nodes
	nodes.Sort()

	tree := BuildHuffmanTree(nodes)

	//pre order traversal
	PreOrderTraversal(tree)
}
