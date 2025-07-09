package table


import "strings"


type Generator interface{
	NewTable(text string) EncodingTable
}

type EncodingTable map[rune]string

type decodingTree struct {
	Data string
	Left *decodingTree
	Right *decodingTree	
}


func (et EncodingTable) Decode(text string)string{
	dt := et.decodingTree()

	return dt.Decode(text)

}

 
func (dt *decodingTree) add(code string, value rune){
	currentNode := dt


	for _, ch := range code{
		switch ch{
			case '0':
				if currentNode.Left == nil{
					currentNode.Left = &decodingTree{}
				}
				currentNode = currentNode.Left 
			case '1':
				if currentNode.Right == nil{
					currentNode.Right = &decodingTree{}
				}
				currentNode = currentNode.Right
 
		}
	}
	currentNode.Data = string(value)

}



func (ec EncodingTable)decodingTree() decodingTree{
	res := decodingTree{}
		
	for ch, code := range ec{
		res.add(code, ch)

	}


	return res
}

func (dt *decodingTree) Decode(bStr string) string{
	var buf strings.Builder

	currentNode := dt
	
	for _, ch := range bStr{
		
		if currentNode.Data != ""{
			buf.WriteString(currentNode.Data)
			currentNode = dt
		}
		
		switch ch{
			case '0':
				currentNode = currentNode.Left
			case '1':
				currentNode = currentNode.Right
		}

	}	
	if currentNode.Data != ""{
			buf.WriteString(currentNode.Data)
			currentNode = dt
	}

	return buf.String()
}
