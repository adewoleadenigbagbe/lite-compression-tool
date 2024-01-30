package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

type HuffmanNodeArray []*HuffmanNode

func (hf HuffmanNodeArray) Len() int { return len(hf) }
func (hf HuffmanNodeArray) Less(i, j int) bool {
	if hf[i].weight == hf[j].weight {
		return hf[i].char < hf[j].char
	}
	return hf[i].weight < hf[j].weight
}
func (hf HuffmanNodeArray) Swap(i, j int) { hf[i], hf[j] = hf[j], hf[i] }
func (hf HuffmanNodeArray) Sort()         { sort.Sort(hf) }

type HuffmanNode struct {
	left   *HuffmanNode
	right  *HuffmanNode
	char   rune
	weight int64
}

func getFile(filename string) (*os.File, error) {
	return os.Open(filename)
}

func writeToFile(filename string, rs []rune) error {
	file, err := os.Create(filename) // Creates or truncates with os.O_RDWR mode
	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, r := range rs {
		_, err = writer.WriteRune(r)
		if err != nil {
			return err
		}
	}

	writer.Flush() // Make sure to flush!

	return nil
}

func calculateCharacterFrequency(file *os.File) (map[rune]int64, []rune, error) {
	var runes []rune
	frequencyTable := make(map[rune]int64)
	reader := bufio.NewReader(file)
	for {
		r, _, err := reader.ReadRune()
		if err != nil && err != io.EOF {
			return nil, nil, err
		}

		if err == io.EOF {
			break
		}

		frequencyTable[r]++
		runes = append(runes, r)
	}

	return frequencyTable, runes, nil
}

func buildHuffmanTree(nodes HuffmanNodeArray) *HuffmanNode {
	for len(nodes) > 1 {
		current := &HuffmanNode{
			weight: nodes[0].weight + nodes[1].weight,
			left:   nodes[0],
			right:  nodes[1],
		}
		nodes = append(nodes, current)
		nodes = nodes[2:]

		nodes.Sort()
	}

	return nodes[0]
}

func generateCodes(root *HuffmanNode, characterCodeMap map[rune]string, frequencyTable map[rune]int64, codes string, code string) {
	if root == nil {
		codes = codes[:len(codes)-1]
		return
	}

	codes += code
	if root.char != 0 {
		characterCodeMap[root.char] = codes
	}

	generateCodes(root.left, characterCodeMap, frequencyTable, codes, "0")
	generateCodes(root.right, characterCodeMap, frequencyTable, codes, "1")
}

func decode(root *HuffmanNode, encoded string) []rune {
	current := root
	var b []rune
	for i := 0; i < len(encoded); i++ {
		if encoded[i] == 48 {
			current = current.left
		} else {
			current = current.right
		}

		if current.left == nil && current.right == nil {
			b = append(b, current.char)
			current = root
		}
	}

	return b
}

func main() {
	var (
		err            error
		file           *os.File
		frequencyTable = make(map[rune]int64)
		characterArray = make([]rune, 0)
	)

	file, err = getFile("lesmiserables.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	frequencyTable, characterArray, err = calculateCharacterFrequency(file)
	if err != nil {
		log.Fatal(err)
	}
	var nodes HuffmanNodeArray
	for k, v := range frequencyTable {
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
	tree := buildHuffmanTree(nodes)

	codes := ""
	m := make(map[rune]string)

	//pre order traversal
	generateCodes(tree, m, frequencyTable, codes, "")

	var sb strings.Builder
	for _, v := range characterArray {
		sb.WriteString(m[v])
	}

	decodedString := decode(tree, sb.String())
	err = writeToFile("decoded_lesmiserables.txt", decodedString)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("done")
}
