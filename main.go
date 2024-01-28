package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type HuffmanNodeArray []*HuffmanNode

func (hf HuffmanNodeArray) Len() int           { return len(hf) }
func (hf HuffmanNodeArray) Less(i, j int) bool { return hf[i].weight < hf[j].weight }
func (hf HuffmanNodeArray) Swap(i, j int)      { hf[i], hf[j] = hf[j], hf[i] }
func (hf HuffmanNodeArray) Sort()              { sort.Sort(hf) }

type HuffmanNode struct {
	left   *HuffmanNode
	right  *HuffmanNode
	char   byte
	weight int64
}

type HuffmanRow struct {
	char      string
	frequency int64
	code      string
	bits      int
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

func BuildHuffmanTree(nodes HuffmanNodeArray) *HuffmanNode {
	current := &HuffmanNode{}
	for len(nodes) >= 2 {
		current.left = nodes[0]
		current.right = nodes[1]
		current.weight = current.left.weight + current.right.weight

		nodes = nodes[1:]
		nodes = append(nodes, current)
		nodes.Sort()
	}

	return current
}

func PreOrderTraversal(root *HuffmanNode, table *[]HuffmanRow, frequencyTable map[byte]int64, codes string, code string) {
	if root == nil {
		codes = codes[:len(codes)-1]
		return
	}

	codes += code
	if root.char != 0 {
		row := HuffmanRow{
			char:      string(root.char),
			frequency: frequencyTable[root.char],
			code:      codes,
			bits:      len(codes),
		}
		*table = append(*table, row)
	}

	PreOrderTraversal(root.left, table, frequencyTable, codes, "0")
	PreOrderTraversal(root.right, table, frequencyTable, codes, "1")
}

func main() {
	// var (
	// 	err            error
	// 	file           *os.File
	// 	frequencyTable = make(map[byte]int64)
	// )

	// file, err = getFile("lesmiserables.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer file.Close()

	// frequencyTable, err = calculateCharacterFrequency(file)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	freqTable := map[byte]int64{
		67: 32,
		68: 42,
		69: 120,
		75: 7,
		76: 42,
		77: 24,
		85: 37,
		90: 2,
	}

	var nodes HuffmanNodeArray
	for k, v := range freqTable {
		node := &HuffmanNode{
			char:   k,
			weight: v,
		}
		nodes = append(nodes, node)
	}

	if len(nodes) < 2 {
		fmt.Println("Huffman tree cannot be built with less than two nodes")
		return
	}

	//sort the nodes
	nodes.Sort()

	// build tree
	tree := BuildHuffmanTree(nodes)

	//pre order traversal
	codes := ""
	var table *[]HuffmanRow = &[]HuffmanRow{}
	PreOrderTraversal(tree, table, freqTable, codes, "")

	fmt.Println("\nPrinting huffman table")

	fmt.Println(table)
}
