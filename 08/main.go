package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	Start       int
	End         int
	NumChildren int
	NumMetadata int
	Children    []*Node
	Metadata    []int
}

var data = make([]int, 0)

func main() {

	filePath := flag.String("file", "input.txt", "file containing the input data")
	flag.Parse()

	lines, err := readFile(*filePath)
	if err != nil {
		log.Fatalf("cannot read file %s: %v", *filePath, err)
	}

	line := lines[0]

	dataStr := strings.Split(line, " ")

	for _, s := range dataStr {
		data = append(data, mustParseInt(s))
	}

	root := new(Node)

	root.Populate()

	var sum int
	root.SumMetadata(&sum)
	fmt.Printf("part 1: %d\n", sum)

	value := root.Value()
	fmt.Printf("part 2: %d\n", value)

	//root.Dump()

}

func (n *Node) Populate() {

	n.NumChildren = data[n.Start]
	n.NumMetadata = data[n.Start+1]

	if n.NumMetadata < 1 {
		log.Fatalf("NumMetadata cannot be < 1: %v\n", n)
	}

	childStart := n.Start + 2
	for len(n.Children) < n.NumChildren {
		child := new(Node)
		child.Start = childStart
		child.Populate()

		n.Children = append(n.Children, child)
		childStart = child.End
	}

	if n.NumChildren == 0 {
		n.End = n.Start + 2 + n.NumMetadata
		n.Metadata = data[n.Start+2 : n.End]
		return
	}

	if num := len(n.Children); num > 0 {
		lastChild := n.Children[num-1]
		if lastChild.End > 0 {
			n.End = lastChild.End + n.NumMetadata
			n.Metadata = data[lastChild.End:n.End]
		}
	}
}

func (n *Node) SumMetadata(sum *int) {

	for _, val := range n.Metadata {
		*sum += val
	}

	for _, child := range n.Children {
		child.SumMetadata(sum)
	}

}

func (n *Node) Value() int {

	// If a node has no child nodes, its value is the sum of its metadata entries.
	// So, the value of node B is 10+11+12=33, and the value of node D is 99.
	if n.NumChildren == 0 {
		var sum int
		for _, val := range n.Metadata {
			sum += val
		}
		return sum
	}

	// However, if a node does have child nodes, the metadata entries become indexes which
	// refer to those child nodes. A metadata entry of 1 refers to the first child node,
	// 2 to the second, 3 to the third, and so on. The value of this node is the sum of the
	// values of the child nodes referenced by the metadata entries. If a referenced child
	// node does not exist, that reference is skipped. A child node can be referenced multiple
	// time and counts each time it is referenced. A metadata entry of 0 does not refer to any child node.

	var sum int
	for _, val := range n.Metadata {
		if val < 1 || val > n.NumChildren {
			continue
		}
		child := n.Children[val-1]
		sum += child.Value()
	}

	return sum

}

func (n *Node) Dump() {
	jsonBytes, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		log.Fatalf("could not dump node as json: %+v", err)
	}
	fmt.Println(string(jsonBytes))
}

func readFile(path string) ([]string, error) {
	if path == "" {
		return nil, errors.New("file path not specified")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func mustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("cannot convert string %s to integer: %v", s, err)
	}
	return i
}
