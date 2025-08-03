package haffman

import (
	"fmt"
	"strings"

	"container/heap"

	"archiver/lib/compression/vlc/table"
)

type Generator struct{}

type charStat map[rune]int

type Node struct{
	Char rune
	Quantite int
	Bits uint32
	Size int
	Left* Node
	Right* Node
}

type Queue []*Node

func (q Queue) Len() int{
	return len(q)
}

func (q Queue) Less(i, j int) bool {
    if q[i].Quantite == q[j].Quantite {
        iLeaf := q[i].Left == nil && q[i].Right == nil
        jLeaf := q[j].Left == nil && q[j].Right == nil
        
        if iLeaf && !jLeaf {
            return false
        }
        if !iLeaf && jLeaf {
            return true
        }
        
        return q[i].Char < q[j].Char
    }
    return q[i].Quantite < q[j].Quantite
}

func (q Queue) Swap(i, j int){
	q[i], q[j] = q[j], q[i]
}

func (q *Queue) Push(x interface{}){
	*q = append(*q, x.(*Node))	
}

func (q *Queue) Pop() interface{}{
	old := *q
	n := len(old)

	item := old[n-1]
	*q = old[0: n-1]

	return item 
}



type encodingTable map[rune]Node

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
	queue := &Queue{}
	heap.Init(queue)

	for ch, qty := range stat{
		item := &Node{
			Char: ch,
			Quantite: qty,
		}
		heap.Push(queue, item)
	}


	if queue.Len() == 0{
		return encodingTable{}
	}

	root := buildHaffmanTree(queue)	

	res := make(encodingTable)
	assignCodes(root, res)	


	return res
}



func buildHaffmanTree(queue *Queue) *Node{
	
	if queue.Len() == 0{
		return nil
	}
		

	for queue.Len() > 1{
		left := heap.Pop(queue).(*Node)
		right := heap.Pop(queue).(*Node)

		parent := &Node{
			Quantite: left.Quantite + right.Quantite,
			Left: left,
			Right: right,
		}

		heap.Push(queue, parent)
	}

	root := heap.Pop(queue).(*Node)


	return root
}


func assignCodes(node *Node, table encodingTable) {
    if node == nil {
        return
    }
    
    var traverse func(n *Node, code uint32, size int)
    traverse = func(n *Node, code uint32, size int) {
        if n == nil {
            return
        }
        
        if n.Left == nil && n.Right == nil {
            table[n.Char] = Node{
                Char:     n.Char,
                Quantite: n.Quantite,
                Bits:     code,
                Size:     size,
            }
            return
        }
        
        if n.Left != nil {
            traverse(n.Left, code<<1, size+1)
        }
        if n.Right != nil {
            traverse(n.Right, (code<<1)|1, size+1)
        }
    }
    
    if node.Left == nil && node.Right == nil {
        table[node.Char] = Node{
            Char:     node.Char,
            Quantite: node.Quantite,
            Bits:     0,
            Size:     1,
        }
    } else {
        traverse(node, 0, 0)
    }
}

func newCharStat(text string) charStat{

	res := make(charStat)

	for _, ch := range text{
		res[ch]++
	}
	
	return res
}
