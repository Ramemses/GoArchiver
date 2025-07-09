package haffman

import (
	"sort"
	"math"
	"fmt"
	"strings"

	"container/heap"

	"archiver/lib/compression/vlc/table"
)

type Generator struct{}

type charStat map[rune]int

type code struct{
	Char rune
	Quantity int
	Bits uint32
	Size int
}	

type encodingTable map[rune]code


// HaffmanTree and Queue

type Node struct{
	Code code
	Left* Node
	Right* Node
}

type pQueue *[]Node

func (pq* pQueue) Len() int {return len(pq)}


func (pq* pQueue) Less(i, j int) bool{
	return pq[i].Code.Quantity < pq[j].Code.Quantity
}

func (pq* pQueue) Swap(i, j int){
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq* pQueue) Push(x interface{}){
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq* pQueue) Pop() interface{}{
	old := *pq
	n := len(old)
	item := old[n-1]
	old = old[0:n-1]

	return item
}

func BuildHaffmanTree(codes []code) Item{
	
	//Here we make queue and fill it
	pq := make(pQueue, 0)
	heap.Init(&pq)

	for code := range codes{
		heap.Push(&pq, &Node{Code: code})
	}

	
	for pq.Len() > 1{
		left := heap.Pop(&pq)
		right := heap.Pop(&pq)
	
		item := Item{Code:}	
	}
	
}



//----------






func NewGenerator() Generator{
	return Generator{}
}


func (g Generator) NewTable(text string) table.EncodingTable{
			
		encTable := build(text)	
		return table.EncodingTable(encTable.Export())
}

func (et encodingTable) Export() map[rune]string{
	res := make(map[rune]string)

	for k, v := range et{
		byteStr := fmt.Sprintf("%b", v.Bits)

		lenDiff := v.Size - len(byteStr)

		if lenDiff > 0{
			byteStr = strings.Repeat("0", lenDiff) + byteStr
		}
		res[k] = byteStr
	}

	return res
}


func build(str string) encodingTable{
	
	stat := newCharStat(str)
	codes := make([]code, 0, len(stat))

	for ch, qty := range stat{
		codes = append(codes, code{
				Char: ch,
				Quantity: qty,
			})
	}

	sort.Slice(codes, func(i, j int) bool {
		if codes[i].Quantity != codes[j].Quantity{
			return codes[i].Quantity > codes[j].Quantity
		}
		return codes[i].Char < codes[j].Char	
	})

//TODO: Make haffman assignCodes
	assignCodes(codes)

	res := make(encodingTable)
	
	for _, code := range codes{
		res[code.Char] = code
	}
	
	return res
}

//TODO: Make Haffman assignCodes
func assignCodes(codes []code){
	if len(codes) == 0{
		return 
	}
	if len(codes) == 1{
		if codes[0].Size == 0{
			codes[0].Bits <<= 1
			codes[0].Size++

		}
		return 
	}

// REPEAT UNTIL Len(queue) > 1

	//Pop 2 elements with least quantites
	
	//Make Node, where fst element equals Node.Left and snd equals Node.Right

	//Insert Node in queue 
	
// Work with HaffmanTree

	//Walk for each Node and:
		// When we go to right, we add '1' into char-code
		// When we go to left, we add '0' into char-code
		// Then we check existing of Node.Code.Char in codes and give codes[Node.Code.Char] the current code if that true

	

	//return codes


}


func abs(x int)int{
	if x < 0{
		x = -x
	}
	return x
}


func newCharStat(text string) charStat{

	res := make(charStat)

	for _, ch := range text{
		res[ch]++
	}
	
	return res
}

